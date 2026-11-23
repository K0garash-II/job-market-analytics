package collect

import (
	"context"
	"errors"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"testing"
)

type feed []models.RawVacancy

func (f feed) Fetch(context.Context, int) ([]models.RawVacancy, error) { return f, nil }

type history map[string]bool

func (h history) Seen(_ context.Context, _ string, id string) (bool, error) { return h[id], nil }
func (h history) Mark(_ context.Context, _ string, id string) error         { h[id] = true; return nil }

type publishFunc func(context.Context, models.RawVacancy) error

func (p publishFunc) Publish(ctx context.Context, v models.RawVacancy) error { return p(ctx, v) }
func TestPublicationFailureRemainsRetryable(t *testing.T) {
	ctx := context.Background()
	seen := history{"old": true}
	jobs := feed{{ExternalID: "old"}, {ExternalID: "new"}}
	calls := 0
	publish := publishFunc(func(_ context.Context, v models.RawVacancy) error {
		if v.ExternalID != "new" {
			t.Fatalf("republished %q", v.ExternalID)
		}
		calls++
		if calls == 1 {
			return errors.New("queue unavailable")
		}
		return nil
	})
	if err := Page(ctx, jobs, 0, seen, publish); err == nil || seen["new"] {
		t.Fatal("failed publication must remain pending")
	}
	for i := 0; i < 2; i++ {
		if err := Page(ctx, jobs, 0, seen, publish); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 2 || !seen["new"] {
		t.Fatalf("calls=%d, history=%v", calls, seen)
	}
}
