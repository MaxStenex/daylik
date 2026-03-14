package httputil

import (
	"encoding/json"
	"net/http"

	api "github.com/maximrozinkevich/daylik/internal/generated/api"
)

const maxBodySize = 1 << 20 // 1 MB

func WriteJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		data = []byte(`{"error":"internal server error"}`)
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func DecodeJSON(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}

func BindJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := DecodeJSON(r, v); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrResp("invalid request body"))
		return false
	}
	return true
}

func ErrResp(msg string) api.ErrorResponse {
	return api.ErrorResponse{Error: msg}
}
