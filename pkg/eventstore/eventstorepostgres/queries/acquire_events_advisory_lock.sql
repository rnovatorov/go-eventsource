SELECT
    pg_advisory_xact_lock('events'::REGCLASS::BIGINT);
