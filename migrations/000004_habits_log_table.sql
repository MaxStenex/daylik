-- +goose Up
CREATE TABLE IF NOT EXISTS habits_log (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id        UUID        NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    unit            TEXT        NOT NULL,
    daily_target    INTEGER     NOT NULL,
    completed_count INTEGER     NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON habits_log(habit_id);
CREATE INDEX ON habits_log(user_id);

-- +goose Down
DROP TABLE IF EXISTS habits_log;
