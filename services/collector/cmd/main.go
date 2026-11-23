package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/collect"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/dedup"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/parser"
	"github.com/malayegor9-dotcom/job-market-analytics/services/collector/internal/queue"
	"log"
	"os"
)

func run(ctx context.Context) error {
	client, err := db.NewRedisClient(ctx, os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}
	defer client.Close()
	nc, err := db.NewNATSConn(os.Getenv("NATS_URL"))
	if err != nil {
		return err
	}
	defer nc.Close()
	js, _, err := db.VacancyStream(ctx, nc)
	if err != nil {
		return err
	}
	history := dedup.New(client)
	publisher := queue.Publisher{JS: js}
	var failures []error
	for _, feed := range []struct {
		source parser.Parser
		pages  int
	}{{parser.NewRemoteOKParser(), 1}, {parser.NewHabrParser(), 3}} {
		for page := 0; page < feed.pages; page++ {
			if err := collect.Page(ctx, feed.source, page, history, publisher); err != nil {
				failures = append(failures, fmt.Errorf("%s page %d: %w", feed.source.Source(), page, err))
			}
		}
	}
	return errors.Join(failures...)
}
func main() {
	_ = godotenv.Load()
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
