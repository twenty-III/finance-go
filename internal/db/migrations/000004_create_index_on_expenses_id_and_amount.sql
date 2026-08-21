-- +goose Up
CREATE INDEX IF NOT EXISTS idx_expenses_id_amount ON expenses(id, amount);

-- +goose Down
DROP INDEX IF EXISTS idx_expenses_id_amount;
