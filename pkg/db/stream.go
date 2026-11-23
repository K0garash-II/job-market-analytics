package db

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const VacancySubject = "vacancies.raw"

func VacancyStream(ctx context.Context, nc *nats.Conn) (jetstream.JetStream, jetstream.Stream, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: "VACANCIES", Subjects: []string{VacancySubject}, Retention: jetstream.WorkQueuePolicy})
	return js, stream, err
}
