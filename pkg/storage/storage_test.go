package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"os"
	"testing"
)

func TestAtomicSaveAndRedelivery(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run with PostgreSQL")
	}
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	// Temporary tables shadow the application's tables for this single connection.
	cfg.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, err = pool.Exec(ctx, `CREATE TEMP TABLE sources(id serial primary key,name text);INSERT INTO sources(name) VALUES ('habr');
 CREATE TEMP TABLE companies(id serial primary key,name text unique);
 CREATE TEMP TABLE vacancies(id serial primary key, external_id text,source_id int,company_id int,title text,description text,location text,salary_min bigint,salary_max bigint,currency text,remote bool,url text,published_at timestamptz,unique(source_id,external_id));
 CREATE TEMP TABLE vacancy_skills(vacancy_id int,skill_id int CHECK(skill_id>0));`)
	if err != nil {
		t.Fatal(err)
	}
	repo := New(pool)
	job := models.RawVacancy{Source: "habr", ExternalID: "1", Title: "Go", Company: "Company"}
	if err := repo.Save(ctx, job, []int{-1}); err == nil {
		t.Fatal("expected skill insert to fail")
	}
	for i := 0; i < 2; i++ {
		if err := repo.Save(ctx, job, []int{1}); err != nil {
			t.Fatal(err)
		}
	}
	var jobs, skills int
	if err := pool.QueryRow(ctx, "SELECT (SELECT count(*) FROM vacancies),(SELECT count(*) FROM vacancy_skills)").Scan(&jobs, &skills); err != nil {
		t.Fatal(err)
	}
	if jobs != 1 || skills != 1 {
		t.Fatalf("got %d jobs and %d skills", jobs, skills)
	}
}
