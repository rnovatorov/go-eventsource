WITH subscription AS (
    SELECT
        position
    FROM
        es_subscriptions
    WHERE
        id = @subscription_id
    FOR NO KEY UPDATE
),
new_events AS (
    SELECT
        e.id,
        e.sequence_number
    FROM
        es_events e
        JOIN subscription s ON e.sequence_number > s.position
),
inserted_events AS (
INSERT INTO es_subscription_backlogs (subscription_id, event_id)
    SELECT
        @subscription_id,
        id
    FROM
        new_events
),
newest_event AS (
    SELECT
        sequence_number
    FROM
        new_events
    ORDER BY
        sequence_number DESC
    LIMIT 1)
UPDATE
    es_subscriptions s
SET
    position = e.sequence_number
FROM
    newest_event e
WHERE
    id = @subscription_id;
