package collect_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/collect"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/dedup"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/queue"
)

type singleJob models.RawVacancy

func (j singleJob) Fetch(context.Context, int) ([]models.RawVacancy, error) {
	return []models.RawVacancy{models.RawVacancy(j)}, nil
}

// Run against a disposable migrated database with a processor listening on the test stream.
func TestCollectorToProcessor(t *testing.T) {
	pg, redisURL, natsURL := os.Getenv("TEST_DATABASE_URL"), os.Getenv("TEST_REDIS_URL"), os.Getenv("TEST_NATS_URL")
	if pg == "" || redisURL == "" || natsURL == "" {
		t.Skip("requires disposable PostgreSQL, Redis, JetStream and processor")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := db.NewPostgresPool(ctx, pg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	redis, err := db.NewRedisClient(ctx, redisURL)
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()
	nc, err := db.NewNATSConn(natsURL)
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()
	js, _, err := db.VacancyStream(ctx, nc)
	if err != nil {
		t.Fatal(err)
	}
	job := singleJob{ExternalID: "pipeline-" + time.Now().Format(time.RFC3339Nano), Source: "habr", Title: "Go and Redis", Company: "Pipeline test"}
	history, out := dedup.New(redis), queue.Publisher{JS: js}
	for i := 0; i < 2; i++ {
		if err := collect.Page(ctx, job, 0, history, out); err != nil {
			t.Fatal(err)
		}
	}
	for {
		var count int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM vacancies v JOIN vacancy_skills vs ON vs.vacancy_id=v.id WHERE v.external_id=$1`, job.ExternalID).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count == 2 {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatalf("expected two extracted skills, got %d", count)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
