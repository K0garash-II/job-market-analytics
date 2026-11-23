package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

type Storage struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Storage { return &Storage{pool} }
func (s *Storage) Skills(ctx context.Context) (map[string]int, error) {
	rows, err := s.pool.Query(ctx, "SELECT name,id FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	index := map[string]int{}
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&name, &id); err != nil {
			return nil, err
		}
		index[name] = id
	}
	return index, rows.Err()
}
func (s *Storage) Save(ctx context.Context, v models.RawVacancy, skills []int) error {
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var company, id, source int
		if err := tx.QueryRow(ctx, "SELECT id FROM sources WHERE name=$1", v.Source).Scan(&source); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `INSERT INTO companies(name) VALUES($1) ON CONFLICT(name) DO UPDATE SET name=EXCLUDED.name RETURNING id`, v.Company).Scan(&company); err != nil {
			return err
		}
		var published any
		if !v.PublishedAt.IsZero() {
			published = v.PublishedAt
		}
		err := tx.QueryRow(ctx, `INSERT INTO vacancies(external_id,source_id,company_id,title,description,location,salary_min,salary_max,currency,remote,url,published_at)
 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT(source_id,external_id) DO NOTHING RETURNING id`, v.ExternalID, source, company, v.Title, v.Description, v.Location, v.SalaryFrom, v.SalaryTo, v.Currency, v.Remote, v.URL, published).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO vacancy_skills(vacancy_id,skill_id) SELECT $1,unnest($2::int[])`, id, skills)
		return err
	})
}
