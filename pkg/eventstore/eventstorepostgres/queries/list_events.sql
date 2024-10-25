SELECT
    id,
    aggregate_id,
    aggregate_version,
    timestamp,
    metadata,
    data
FROM
    es_events
WHERE
    aggregate_id = @aggregate_id
ORDER BY
    aggregate_version;
