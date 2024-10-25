UPDATE
    es_aggregates
SET
    version = @new_aggregate_version
WHERE
    id = @aggregate_id
    AND version = @expected_aggregate_version;
