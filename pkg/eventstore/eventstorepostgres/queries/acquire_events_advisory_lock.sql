SELECT
    pg_advisory_xact_lock('es_events'::REGCLASS::BIGINT);
