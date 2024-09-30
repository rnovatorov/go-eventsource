BEGIN;

CREATE SCHEMA IF NOT EXISTS eventsource;

CREATE TABLE IF NOT EXISTS eventsource.aggregates (
    id TEXT PRIMARY KEY,
    version INT NOT NULL CHECK (version >= 0)
);

CREATE TABLE IF NOT EXISTS eventsource.events (
    id TEXT PRIMARY KEY,
    aggregate_id TEXT NOT NULL REFERENCES eventsource.aggregates (id),
    aggregate_version INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL,
    data JSONB NOT NULL,
    UNIQUE (aggregate_id, aggregate_version)
);

END;
