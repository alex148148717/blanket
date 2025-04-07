-- +goose Up
-- +goose StatementBegin
CREATE TABLE property_transactions
(
    id              UInt64,
    user_id         String,
    property_id     String,
    amount          Float64,
    date            DateTime,
    created_at      DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY (date, id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_transactions;
-- +goose StatementEnd
