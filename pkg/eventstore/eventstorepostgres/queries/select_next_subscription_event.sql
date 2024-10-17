WITH subscription AS (
    SELECT
        position
    FROM
        subscriptions
    WHERE
        id = @subscription_id
    FOR NO KEY UPDATE
)
SELECT
    id,
    aggregate_id,
    aggregate_version,
    timestamp,
    metadata,
    data
FROM
    events e
    JOIN subscription s ON e.sequence_number = s.position + 1;
