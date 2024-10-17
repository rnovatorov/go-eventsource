INSERT INTO aggregates (id, version)
    VALUES (@aggregate_id, 0)
ON CONFLICT
    DO NOTHING;
