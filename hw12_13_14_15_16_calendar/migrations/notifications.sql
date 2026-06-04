-- +goose Up
-- +goose StatementBegin
CREATE TABLE notification
(
    event_id   bigint      primary key,
    title      text        not null,
    event_date timestamptz not null,
    user_id    bigint      not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notification;
-- +goose StatementEnd
