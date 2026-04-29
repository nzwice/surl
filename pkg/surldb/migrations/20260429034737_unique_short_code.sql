-- +goose Up
CREATE UNIQUE INDEX idx_short_code_unique ON urls(short_code);
ALTER TABLE urls ADD CONSTRAINT unique_short_code UNIQUE USING INDEX idx_short_code_unique;

-- +goose Down
DROP INDEX IF EXISTS idx_short_code_unique;
