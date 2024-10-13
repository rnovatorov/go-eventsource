SELECT
    id,
    aggregate_id,
    aggregate_version,
    timestamp,
    metadata,
    data
FROM
    events
WHERE
    aggregate_id = @aggregate_id;
