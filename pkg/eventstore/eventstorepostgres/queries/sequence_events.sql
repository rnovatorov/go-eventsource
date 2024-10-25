WITH last_sequenced_event AS (
    SELECT
        coalesce(max(sequence_number), 0) AS sequence_number
    FROM
        es_events
),
non_sequenced_events AS (
    SELECT
        id,
        row_number() OVER (ORDER BY aggregate_version) AS row_number
    FROM
        es_events
    WHERE
        sequence_number IS NULL)
UPDATE
    es_events e
SET
    sequence_number = l.sequence_number + n.row_number
FROM
    non_sequenced_events n
    CROSS JOIN last_sequenced_event l
WHERE
    e.id = n.id;
