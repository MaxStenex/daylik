# Habit Progression Logs — Design

Status: planned, not yet implemented.

## Goal

Persist the user's per-day habit progression so that:

1. A user can open a habit and see a historical calendar of completions.
2. The dashboard can render a GitHub-style contribution grid from aggregated logs.
3. Past days are **immutable** under later habit edits — if the user changes a habit's
   target from 300 to 500 today, yesterday's row still reads 300/300, not 300/500.

## Model

Two tables:

- `habit_daily_progress` — mutable running state for the **current** day.
  One row per habit, updated via `SetProgress`. When the user finalizes the day,
  its contents are copied into `habit_logs` and the row is reset / removed.
- `habit_logs` — immutable, append-only history. One row per habit per day, created
  only by finalization. Each row stores a **snapshot** of the habit's
  `daily_target`, `unit`, and `exp_reward` at the moment of freeze. This is how
  past rows stay correct after habit edits: they carry their own copy of the values.

No per-day versioning on `habits` itself; edits mutate the habit in place and
only affect days that have not yet been finalized.

## Invariants

- `UNIQUE (habit_id, log_date)` on `habit_logs` — finalize is idempotent per day.
- `habit_logs.habit_id` has **no** `ON DELETE CASCADE`. Habit deletion is a soft
  delete; history must survive it.
- `habit_logs` has no `deleted_at`. Rows are never modified or removed after write.
- `habit_daily_progress` has exactly zero or one row per habit; its `log_date`
  tells you which day the stored `progress` applies to. A write for a newer date
  overwrites (resets) the row.
- Progress is set by **absolute value**, not delta. Repeated identical writes are
  no-ops; retries are safe.

## Dates and timezones

For the MVP, `log_date` is computed in **UTC**. This is known to be wrong for a
leaderboard scenario (a Tokyo user's "day" would cut off at 09:00 local), and
should be revisited before the leaderboard ships. The expected migration path is
to add `users.timezone TEXT NOT NULL DEFAULT 'UTC'` and resolve `log_date` per user.
Not done now to keep the first cut small.

## Finalization trigger

Finalization is triggered by an **explicit user action** ("Finish day" button)
hitting a dedicated endpoint. Automatic end-of-day finalization via cron is
**not** implemented in this iteration, but the usecase is shaped so that adding it
later is a matter of calling the same function from a scheduled job for every
user with un-finalized progress — no business logic changes required.

Concretely: the `FinalizeDay` usecase takes a `userID` and a `date`, reads the
relevant `habit_daily_progress` rows, writes `habit_logs`, and clears the
progress rows in one transaction. A future cron worker will iterate users and
call the same function.

## Schema

### Migration `000004_create_habit_logs.sql`

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS habit_logs (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id     UUID        NOT NULL REFERENCES habits(id),
    log_date     DATE        NOT NULL,
    progress     BIGINT      NOT NULL,
    daily_target BIGINT      NOT NULL,
    unit         TEXT        NOT NULL,
    exp_reward   BIGINT      NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (habit_id, log_date)
);

CREATE INDEX ON habit_logs (user_id, log_date);
CREATE INDEX ON habit_logs (habit_id, log_date DESC);

-- +goose Down
DROP TABLE IF EXISTS habit_logs;
```

Note: `user_id` is denormalized onto `habit_logs` so that user-scoped queries
(dashboard grid, weekly recap) don't need to join `habits`. The trade-off is
accepted because logs are append-only and the denormalized field can never
drift.

### Migration `000005_create_habit_daily_progress.sql`

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS habit_daily_progress (
    habit_id   UUID        PRIMARY KEY REFERENCES habits(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_date   DATE        NOT NULL,
    progress   BIGINT      NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON habit_daily_progress (user_id, log_date);

-- +goose Down
DROP TABLE IF EXISTS habit_daily_progress;
```

`ON DELETE CASCADE` on `habit_id` is safe here because this row is mutable and
represents "right now"; losing it when a habit is hard-deleted is fine.
(Today's habits use soft delete, so this cascade only fires if the policy ever
changes.)

## Domain layer

### `internal/domain/habit_log/`

- `entity.go`
  ```go
  type HabitLog struct {
      ID          uuid.UUID
      UserID      uuid.UUID
      HabitID     uuid.UUID
      LogDate     time.Time // date only
      Progress    int64
      DailyTarget int64
      Unit        string
      ExpReward   int64
      CreatedAt   time.Time
  }
  ```
- `repository.go`
  ```go
  type Repository interface {
      Create(ctx context.Context, log *HabitLog) error
      FindByHabitAndDateRange(ctx context.Context, habitID uuid.UUID, from, to time.Time) ([]HabitLog, error)
      FindByUserAndDateRange(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]HabitLog, error)
  }
  ```
- `errors.go` — `ErrAlreadyFinalized` (mapped from the UNIQUE violation).

### `internal/domain/habit_progress/`

- `entity.go`
  ```go
  type HabitProgress struct {
      HabitID   uuid.UUID
      UserID    uuid.UUID
      LogDate   time.Time
      Progress  int64
      UpdatedAt time.Time
  }
  ```
- `repository.go`
  ```go
  type Repository interface {
      Upsert(ctx context.Context, p *HabitProgress) error
      FindByHabitID(ctx context.Context, habitID uuid.UUID) (*HabitProgress, error)
      FindByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]HabitProgress, error)
      DeleteByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) error
  }
  ```

## Usecase layer

### `internal/usecase/habit_log/`

- `SetProgress(ctx, SetProgressInput) error`
  - Input: `UserID`, `HabitID`, `Progress` (absolute value).
  - Loads the habit, checks ownership, checks it isn't soft-deleted.
  - Rejects `Progress < 0`. Does not cap at `DailyTarget` — overshoot is allowed
    and preserved as-is.
  - Upserts `habit_daily_progress` with `log_date = today` (UTC).
  - If the existing row is for a stale date, it is overwritten — this is
    effectively a day rollover for that habit.
  - Errors: `ErrNotFound`, `ErrForbidden`, `ErrInvalidProgress`.
  - Note: this does **not** auto-finalize the previous day's row. That's
    `FinalizeDay`'s job. A stale row being overwritten means the prior day was
    never finalized, which is a UX concern — the client should call
    `FinalizeDay` before bedtime, and the future cron is the safety net.

- `FinalizeDay(ctx, FinalizeDayInput) error`
  - Input: `UserID`, optionally `Date` (defaults to today UTC — the param exists
    specifically so a cron can pass yesterday).
  - Reads all `habit_daily_progress` rows for the user at that date.
  - For each one, loads the habit to capture the current `daily_target`, `unit`,
    `exp_reward`, and writes one `habit_logs` row.
  - Deletes the processed `habit_daily_progress` rows.
  - All in a single DB transaction.
  - If a log row for `(habit_id, log_date)` already exists, the usecase treats
    that habit as already finalized and skips it. Returning `ErrAlreadyFinalized`
    only if **nothing** could be written.
  - Side-effect hooks for EXP / achievements / streaks are out of scope for this
    doc but this is the single place they will eventually attach.

- `ListForHabit(ctx, ListForHabitInput) ([]HabitLog, error)` — calendar drill-down.
- `ListForUser(ctx, ListForUserInput) ([]HabitLog, error)` — dashboard grid data.

Both list usecases validate `from <= to` and cap the range (suggest 366 days).

## Delivery layer

### `internal/delivery/http/habit_log/`

All routes sit behind the existing auth middleware group.

- `PUT /api/v1/habits/{id}/progress`
  - Body: `{ "progress": 150 }`
  - Absolute SetProgress (not a delta; see §Daily Progress in PRODUCT.md).
  - 204 on success; 400 invalid; 403 not owner; 404 not found.

- `POST /api/v1/habits/logs/finalize`
  - Body: none (date is derived server-side as today UTC for the MVP).
  - 204 on success; 409 `ErrAlreadyFinalized` only if fully redundant.

- `GET /api/v1/habits/{id}/logs?from=YYYY-MM-DD&to=YYYY-MM-DD`
  - Returns `[]HabitLog` for a single habit — used by the habit detail calendar.
  - Each returned row carries its own `daily_target` / `unit` / `exp_reward`,
    which is exactly what lets the UI render "300/300" for a past day even after
    the habit was edited to target 500.

- `GET /api/v1/habits/logs?from=YYYY-MM-DD&to=YYYY-MM-DD`
  - Returns `[]HabitLog` across all of the user's habits — used by the
    GitHub-style grid. The frontend groups by date and computes color intensity.

## Wiring

- Register `habit_log.Repository` and `habit_progress.Repository` (both Postgres
  implementations) in `cmd/service/main.go`.
- Build a `habit_log` usecase service wired with both repos and the existing
  `habit.Repository`.
- Construct a `habit_log` HTTP handler and mount the routes in
  `internal/delivery/http/router.go`, inside the authenticated group.

## Habit edit / delete interaction

- **Edit**: `internal/usecase/habit/update.go` does not need changes. Frozen
  `habit_logs` rows are immune because they carry their own snapshot; the
  current day's `habit_daily_progress` row picks up the new target automatically
  on its next write because the habit values are resolved at finalize time.
- **Delete**: already soft-delete. `habit_logs` has no cascade so historical rows
  survive. List endpoints join `habits` without filtering on `deleted_at` so
  that past rows for deleted habits still render in the history UI.

## Out of scope for this iteration

- Cron-driven automatic finalization (designed for, but not implemented).
- EXP / streak / achievement side effects on finalization.
- Per-user timezones.
- Retroactive editing of past logs (explicitly disallowed).
