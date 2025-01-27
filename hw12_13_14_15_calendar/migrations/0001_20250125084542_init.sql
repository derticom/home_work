-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    id            UUID PRIMARY KEY,
    header        VARCHAR(100)             NOT NULL,
    date          TIMESTAMP WITH TIME ZONE NOT NULL,
    duration      BIGINT,
    description   TEXT,
    notify_before BIGINT                   NOT NULL
);

-- +goose Down
drop table events;
