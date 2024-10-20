-- River migration 005 [down]
--
-- Revert to migration table based only on `(version)`.
--
-- If any non-main migrations are present, 005 is considered irreversible.
--
DO $body$
BEGIN
    -- Tolerate users who may be using their own migration system rather than
    -- River's. If they are, they will have skipped version 001 containing
    -- `CREATE TABLE river_migration`, so this table won't exist.
    IF (
        SELECT
            to_regclass ('river_migration') IS NOT NULL) THEN
        IF EXISTS (
            SELECT
                *
            FROM
                river_migration
            WHERE
                line <> 'main') THEN
        RAISE EXCEPTION 'Found non-main migration lines in the database; version 005 migration is irreversible because it would result in loss of migration information.';
    END IF;
    ALTER TABLE river_migration RENAME TO river_migration_old;
    CREATE TABLE river_migration (
        id BIGSERIAL PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW( ),
        version BIGINT NOT NULL,
        CONSTRAINT version CHECK (version >= 1 )
    );
    CREATE UNIQUE INDEX ON river_migration USING btree (version);
    INSERT INTO river_migration (created_at, version)
    SELECT
        created_at,
        version
    FROM
        river_migration_old;
    DROP TABLE river_migration_old;
END IF;
END;
$body$
LANGUAGE 'plpgsql';

--
-- Drop `river_job.unique_key`.
--
ALTER TABLE river_job
    DROP COLUMN unique_key;

--
-- Drop `river_client` and derivative.
--
DROP TABLE river_client_queue;

DROP TABLE river_client;
