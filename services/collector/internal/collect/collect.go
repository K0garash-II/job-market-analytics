package collect

import (
	"context"
	"errors"
	"fmt"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
)

type Source interface {
	Fetch(context.Context, int) ([]models.RawVacancy, error)
}
type History interface {
	Seen(context.Context, string, string) (bool, error)
	Mark(context.Context, string, string) error
}
type Publisher interface {
	Publish(context.Context, models.RawVacancy) error
}

func Page(ctx context.Context, src Source, page int, history History, out Publisher) error {
	jobs, err := src.Fetch(ctx, page)
	if err != nil {
		return err
	}
	var failures []error
	for _, job := range jobs {
		seen, err := history.Seen(ctx, job.Source, job.ExternalID)
		if err == nil && !seen {
			err = out.Publish(ctx, job)
			if err == nil {
				err = history.Mark(ctx, job.Source, job.ExternalID)
			}
		}
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", job.ExternalID, err))
		}
	}
	return errors.Join(failures...)
}
