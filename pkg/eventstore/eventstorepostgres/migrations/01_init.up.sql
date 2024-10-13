BEGIN;

CREATE TABLE aggregates (
    id TEXT PRIMARY KEY,
    version INT NOT NULL
);

CREATE TABLE events (
    id TEXT PRIMARY KEY,
    sequence_number BIGINT UNIQUE,
    aggregate_id TEXT NOT NULL REFERENCES aggregates (id),
    aggregate_version INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL,
    data JSONB NOT NULL,
    UNIQUE (aggregate_id, aggregate_version)
);

CREATE INDEX ON events (aggregate_version) INCLUDE (id)
WHERE
    sequence_number IS NULL;

CREATE TABLE subscriptions (
    id TEXT PRIMARY KEY,
    position BIGINT NOT NULL DEFAULT 0
);

END;
