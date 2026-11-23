package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/market"
)

type Reader struct{ Pool *pgxpool.Pool }
type querier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func many[T any](ctx context.Context, db querier, query string, args ...any) ([]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByPos[T])
}

const fields = `v.id,v.title,coalesce(v.location,''),v.salary_min,v.salary_max,coalesce(v.currency,''),coalesce(v.remote,false),coalesce(v.url,''),v.published_at,coalesce(c.name,''),coalesce(s.name,'')`
const joins = ` FROM vacancies v LEFT JOIN companies c ON c.id=v.company_id LEFT JOIN sources s ON s.id=v.source_id `

func (r Reader) Vacancies(ctx context.Context, f market.Filter) (market.Listing, error) {
	const where = ` WHERE ($1='' OR v.location ILIKE '%'||$1||'%')
   AND ($2::boolean IS NULL OR v.remote=$2)
   AND ($3='' OR s.name=$3)
   AND ($4='' OR v.title ILIKE '%'||$4||'%' OR c.name ILIKE '%'||$4||'%')`
	result := market.Listing{Page: f.Page, Limit: f.Limit}
	args := []any{f.Location, f.Remote, f.Source, f.Query}
	err := pgx.BeginTxFunc(ctx, r.Pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, "SELECT count(*)"+joins+where, args...).Scan(&result.Total); err != nil {
			return err
		}
		query := "SELECT " + fields + ", ''" + joins + where + " ORDER BY v.published_at DESC NULLS LAST, v.id DESC LIMIT $5 OFFSET $6"
		var err error
		result.Data, err = many[market.Vacancy](ctx, tx, query, append(args, f.Limit, int64(f.Page-1)*int64(f.Limit))...)
		return err
	})
	return result, err
}

func (r Reader) Vacancy(ctx context.Context, id int) (market.Vacancy, error) {
	rows, err := r.Pool.Query(ctx, "SELECT "+fields+", coalesce(v.description,'')"+joins+"WHERE v.id=$1", id)
	if err != nil {
		return market.Vacancy{}, err
	}
	v, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[market.Vacancy])
	if errors.Is(err, pgx.ErrNoRows) {
		return v, market.ErrNotFound
	}
	return v, err
}
func (r Reader) Search(ctx context.Context, q string) ([]market.SearchResult, error) {
	return many[market.SearchResult](ctx, r.Pool, `SELECT v.id,v.title,coalesce(v.location,''),coalesce(v.remote,false),coalesce(v.url,''),coalesce(c.name,''),ts_rank(v.search_vector,query) AS rank FROM vacancies v LEFT JOIN companies c ON c.id=v.company_id CROSS JOIN websearch_to_tsquery('russian',$1) query WHERE v.search_vector @@ query ORDER BY rank DESC,v.id DESC LIMIT 20`, q)
}
func (r Reader) Skills(ctx context.Context) ([]market.Skill, error) {
	return many[market.Skill](ctx, r.Pool, `SELECT s.name,coalesce(s.category,''),count(*) AS n FROM vacancy_skills vs JOIN skills s ON s.id=vs.skill_id GROUP BY s.id ORDER BY n DESC,s.name LIMIT 20`)
}
func (r Reader) Companies(ctx context.Context) ([]market.Count, error) {
	return many[market.Count](ctx, r.Pool, `SELECT c.name,count(*) AS n FROM vacancies v JOIN companies c ON c.id=v.company_id GROUP BY c.id ORDER BY n DESC,c.name LIMIT 10`)
}
func (r Reader) Salary(ctx context.Context) (market.Salary, error) {
	var s market.Salary
	err := r.Pool.QueryRow(ctx, `SELECT count(*),count(*) FILTER(WHERE salary_min IS NOT NULL OR salary_max IS NOT NULL),coalesce(avg(salary_min),0),coalesce(avg(salary_max),0) FROM vacancies`).Scan(&s.Total, &s.WithSalary, &s.AvgMin, &s.AvgMax)
	return s, err
}
func (r Reader) Overview(ctx context.Context, latestPublication bool) (market.Overview, error) {
	var o market.Overview
	err := pgx.BeginTxFunc(ctx, r.Pool, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `SELECT count(*),count(*) FILTER(WHERE remote),count(DISTINCT company_id),count(*) FILTER(WHERE salary_min IS NOT NULL OR salary_max IS NOT NULL) FROM vacancies`).Scan(&o.Total, &o.Remote, &o.Companies, &o.WithSalary)
		if err != nil {
			return err
		}
		o.Sources, err = many[market.Count](ctx, tx, `SELECT coalesce(s.name,'—'),count(*) AS n FROM vacancies v LEFT JOIN sources s ON s.id=v.source_id GROUP BY s.name ORDER BY n DESC,s.name`)
		if err != nil {
			return err
		}
		o.Days, err = many[market.Count](ctx, tx, `
			WITH anchor AS (
				SELECT CASE WHEN $1::boolean THEN COALESCE(
					(SELECT max(published_at)::date FROM vacancies WHERE published_at < current_date+1), current_date
				) ELSE current_date END AS last_day
			), days AS (SELECT last_day-i AS day FROM anchor, generate_series(0,29) i)
			SELECT to_char(d.day,'YYYY-MM-DD'), count(v.id)
			FROM days d LEFT JOIN vacancies v ON v.published_at >= d.day AND v.published_at < d.day+1
			GROUP BY d.day ORDER BY d.day`, latestPublication)
		return err
	})
	return o, err
}
