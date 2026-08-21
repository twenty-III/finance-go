-- +goose Up
CREATE INDEX IF NOT EXISTS idx_expenses_id_created_at ON expenses (id, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_expenses_id_created_at;
