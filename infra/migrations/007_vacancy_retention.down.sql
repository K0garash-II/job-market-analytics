DROP TRIGGER IF EXISTS vacancies_skip_expired_insert ON vacancies;
DROP FUNCTION IF EXISTS prune_expired_vacancies();
DROP FUNCTION IF EXISTS skip_expired_vacancy();
DROP FUNCTION IF EXISTS vacancy_expiration_cutoff();
