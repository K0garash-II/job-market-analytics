-- Run in a migrated test database. All writes and table substitutions are session-local.
\set ON_ERROR_STOP on
BEGIN;
CREATE TEMP TABLE vacancies (id int PRIMARY KEY, published_at timestamptz);
CREATE TEMP TABLE vacancy_skills (
    vacancy_id int REFERENCES vacancies(id) ON DELETE CASCADE
);
INSERT INTO vacancies VALUES
    (1, vacancy_expiration_cutoff() - INTERVAL '1 microsecond'),
    (2, vacancy_expiration_cutoff()),
    (3, vacancy_expiration_cutoff() + INTERVAL '1 microsecond'),
    (4, NULL);
INSERT INTO vacancy_skills VALUES (1), (2);
CREATE TRIGGER vacancies_skip_expired_insert BEFORE INSERT ON vacancies
FOR EACH ROW EXECUTE FUNCTION skip_expired_vacancy();
INSERT INTO vacancies VALUES (5, now() - INTERVAL '31 days');

DO $$
BEGIN
    IF prune_expired_vacancies() <> 1 THEN RAISE EXCEPTION 'expected one expired vacancy'; END IF;
    IF (SELECT array_agg(id ORDER BY id) FROM vacancies) <> ARRAY[2,3,4] THEN
        RAISE EXCEPTION 'incorrect retention boundary or expired reimport';
    END IF;
    IF (SELECT array_agg(vacancy_id) FROM vacancy_skills) <> ARRAY[2] THEN
        RAISE EXCEPTION 'incorrect cascade';
    END IF;
    IF prune_expired_vacancies() <> 0 THEN RAISE EXCEPTION 'cleanup must be idempotent'; END IF;
END;
$$;
ROLLBACK;
