WITH processible_events AS (
    SELECT DISTINCT ON (e.aggregate_id)
        e.id,
        e.aggregate_id,
        e.aggregate_version,
        e.timestamp,
        e.metadata,
        e.data
    FROM
        es_subscription_backlogs b
        JOIN es_events e ON b.event_id = e.id
    WHERE
        b.subscription_id = @subscription_id
    ORDER BY
        e.aggregate_id,
        e.aggregate_version
)
SELECT
    p.*
FROM
    processible_events p
    JOIN es_subscription_backlogs b ON p.id = b.event_id
WHERE
    b.subscription_id = @subscription_id
LIMIT 1
FOR UPDATE
    OF b SKIP LOCKED;
