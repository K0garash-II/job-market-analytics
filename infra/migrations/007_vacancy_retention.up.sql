CREATE OR REPLACE FUNCTION vacancy_expiration_cutoff()
RETURNS timestamptz LANGUAGE sql STABLE AS $$
    SELECT CURRENT_TIMESTAMP - INTERVAL '720 hours';
$$;

CREATE OR REPLACE FUNCTION skip_expired_vacancy()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.published_at < vacancy_expiration_cutoff() THEN
        RETURN NULL;
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER vacancies_skip_expired_insert
BEFORE INSERT ON vacancies
FOR EACH ROW EXECUTE FUNCTION skip_expired_vacancy();

CREATE OR REPLACE FUNCTION prune_expired_vacancies()
RETURNS bigint LANGUAGE sql AS $$
    WITH expired AS (
        DELETE FROM vacancies
        WHERE published_at < vacancy_expiration_cutoff()
        RETURNING 1
    )
    SELECT count(*) FROM expired;
$$;
