package queue

import (
	"context"
	"encoding/json"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct{ JS jetstream.JetStream }

func (p Publisher) Publish(ctx context.Context, v models.RawVacancy) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = p.JS.Publish(ctx, db.VacancySubject, data, jetstream.WithMsgID(v.Source+":"+v.ExternalID))
	return err
}
