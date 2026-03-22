-- +goose Up
CREATE TABLE IF NOT EXISTS habits (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    exp_reward   BIGINT      NOT NULL,
    daily_target BIGINT      NOT NULL,
    unit         TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    archived_at  TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX ON habits(user_id);

-- +goose Down
DROP TABLE IF EXISTS habits;
