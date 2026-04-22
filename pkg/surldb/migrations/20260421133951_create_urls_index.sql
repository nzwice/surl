-- +goose Up
CREATE INDEX idx_short_code ON urls USING HASH (short_code);
CREATE INDEX idx_expired_at ON urls USING BTREE (expired_at);

-- +goose Down
DROP INDEX idx_short_code;
DROP INDEX idx_expired_at;
