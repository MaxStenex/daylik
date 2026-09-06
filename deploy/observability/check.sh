#!/usr/bin/env bash

set -uo pipefail

COLLECTOR_HEALTH=${COLLECTOR_HEALTH:-http://localhost:13133}
COLLECTOR_OTLP=${COLLECTOR_OTLP:-http://localhost:4318}
JAEGER=${JAEGER:-http://localhost:16686}
LOKI=${LOKI:-http://localhost:3100}
PROMETHEUS=${PROMETHEUS:-http://localhost:9090}
GRAFANA=${GRAFANA:-http://localhost:${GRAFANA_PORT:-3001}}
SERVICE_NAME=${SERVICE_NAME:-daylik}

pass=0
fail=0

green() { printf '\033[32m%s\033[0m\n' "$1"; }
red()   { printf '\033[31m%s\033[0m\n' "$1"; }
dim()   { printf '\033[2m%s\033[0m\n' "$1"; }

ok()   { green "  PASS  $1"; pass=$((pass+1)); }
bad()  { red   "  FAIL  $1"; fail=$((fail+1)); }

retry() {
  local attempts=$1 delay=$2; shift 2
  local i=1
  while [ "$i" -le "$attempts" ]; do
    if "$@"; then return 0; fi
    sleep "$delay"
    i=$((i+1))
  done
  return 1
}

need() {
  command -v "$1" >/dev/null 2>&1 || { red "missing required tool: $1"; exit 1; }
}
need curl
need openssl
need python3

echo
echo "== reachability (from the host, through published ports) =="

http_ok() { curl -fsS --max-time 3 "$1" >/dev/null 2>&1; }

retry 30 2 http_ok "$COLLECTOR_HEALTH"  && ok "collector health ($COLLECTOR_HEALTH)"  || bad "collector health ($COLLECTOR_HEALTH)"
retry 30 2 http_ok "$LOKI/ready"        && ok "loki ready"                            || bad "loki ready"
retry 30 2 http_ok "$PROMETHEUS/-/ready"&& ok "prometheus ready"                      || bad "prometheus ready"
retry 30 2 http_ok "$GRAFANA/api/health"&& ok "grafana health"                         || bad "grafana health"
retry 30 2 http_ok "$JAEGER/"           && ok "jaeger ui"                              || bad "jaeger ui"

echo
echo "== telemetry round-trip =="

TRACE_ID=$(openssl rand -hex 16)
SPAN_ID=$(openssl rand -hex 8)
NOW_NS=$(python3 -c 'import time; print(int(time.time()*1e9))')
END_NS=$((NOW_NS + 5000000))

dim "  trace_id=$TRACE_ID span_id=$SPAN_ID"

trace_payload() {
  cat <<JSON
{"resourceSpans":[{"resource":{"attributes":[
  {"key":"service.name","value":{"stringValue":"$SERVICE_NAME"}}]},
 "scopeSpans":[{"scope":{"name":"obs-check"},"spans":[{
   "traceId":"$TRACE_ID","spanId":"$SPAN_ID","name":"obs-check",
   "kind":2,"startTimeUnixNano":"$NOW_NS","endTimeUnixNano":"$END_NS",
   "status":{}}]}]}]}
JSON
}

log_payload() {
  cat <<JSON
{"resourceLogs":[{"resource":{"attributes":[
  {"key":"service.name","value":{"stringValue":"$SERVICE_NAME"}}]},
 "scopeLogs":[{"scope":{"name":"obs-check"},"logRecords":[{
   "timeUnixNano":"$NOW_NS","severityNumber":9,"severityText":"INFO",
   "body":{"stringValue":"obs-check synthetic log $TRACE_ID"},
   "traceId":"$TRACE_ID","spanId":"$SPAN_ID"}]}]}]}
JSON
}

post_otlp() {
  local path=$1 body=$2
  curl -fsS --max-time 5 -X POST "$COLLECTOR_OTLP$path" \
    -H 'Content-Type: application/json' --data "$body" >/dev/null 2>&1
}

post_otlp /v1/traces "$(trace_payload)" && ok "pushed OTLP trace" || bad "pushed OTLP trace"
post_otlp /v1/logs   "$(log_payload)"   && ok "pushed OTLP log"   || bad "pushed OTLP log"

trace_landed() {
  curl -fsS --max-time 5 "$JAEGER/api/traces/$TRACE_ID" 2>/dev/null \
    | python3 -c 'import json,sys; d=json.load(sys.stdin); sys.exit(0 if d.get("data") else 1)' 2>/dev/null
}
retry 20 2 trace_landed && ok "trace visible in jaeger" || bad "trace visible in jaeger"

log_landed() {
  local start=$((NOW_NS - 60000000000)) end=$((NOW_NS + 60000000000))
  curl -fsS --max-time 5 -G "$LOKI/loki/api/v1/query_range" \
    --data-urlencode "query={service_name=\"$SERVICE_NAME\"}" \
    --data-urlencode "start=$start" --data-urlencode "end=$end" 2>/dev/null \
    | grep -q "$TRACE_ID"
}
retry 20 2 log_landed && ok "log visible in loki" || bad "log visible in loki"

echo
echo "== cross-service wiring =="

targets_up() {
  curl -fsS --max-time 5 "$PROMETHEUS/api/v1/targets" 2>/dev/null | python3 -c '
import json,sys
d=json.load(sys.stdin)
act=d.get("data",{}).get("activeTargets",[])
down=[t["labels"].get("job") for t in act if t.get("health")!="up"]
if not act or down:
    print("down:", ", ".join(sorted(set(down))) or "none active", file=sys.stderr)
    sys.exit(1)
' 2>/dev/null
}
retry 20 3 targets_up && ok "all prometheus targets up" || bad "all prometheus targets up"

ds_health() {
  local uid=$1
  curl -fsS --max-time 8 "$GRAFANA/api/datasources/uid/$uid/health" 2>/dev/null \
    | python3 -c 'import json,sys; sys.exit(0 if json.load(sys.stdin).get("status")=="OK" else 1)' 2>/dev/null
}
for uid in daylik-prometheus daylik-loki daylik-jaeger; do
  retry 10 2 ds_health "$uid" && ok "grafana datasource $uid" || bad "grafana datasource $uid"
done

echo
if [ "$fail" -eq 0 ]; then
  green "all $pass checks passed"
  echo
  dim "  grafana    $GRAFANA"
  dim "  jaeger     $JAEGER/trace/$TRACE_ID"
  dim "  prometheus $PROMETHEUS"
  exit 0
fi
red "$fail of $((pass+fail)) checks failed"
exit 1
