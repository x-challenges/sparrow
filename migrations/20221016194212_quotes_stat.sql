-- +goose Up
-- +goose StatementBegin
CREATE TABLE quotes_stat (
  id              string NOT NULL,
  quotes          json NOT NULL,
  created_at      timestamp NOT NULL,

  PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quotes_stat;
-- +goose StatementEnd
