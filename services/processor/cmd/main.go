package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/models"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/storage"
	"github.com/malayegor9-dotcom/job-market-analytics/services/processor/internal/extractor"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func run(ctx context.Context) error {
	pool, err := db.NewPostgresPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()
	store := storage.New(pool)
	index, err := store.Skills(ctx)
	if err != nil {
		return err
	}
	extract := extractor.New(index)
	nc, err := db.NewNATSConn(os.Getenv("NATS_URL"))
	if err != nil {
		return err
	}
	defer nc.Close()
	_, stream, err := db.VacancyStream(ctx, nc)
	if err != nil {
		return err
	}
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{Durable: "processor", AckPolicy: jetstream.AckExplicitPolicy})
	if err != nil {
		return err
	}
	sub, err := consumer.Consume(func(msg jetstream.Msg) {
		var job models.RawVacancy
		if err := json.Unmarshal(msg.Data(), &job); err != nil {
			log.Print(err)
			_ = msg.Term()
			return
		}
		if err := store.Save(ctx, job, extract.Extract(job.Title, job.Description)); err != nil {
			log.Print(err)
			_ = msg.NakWithDelay(5 * time.Second)
			return
		}
		if err := msg.Ack(); err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		return err
	}
	<-ctx.Done()
	sub.Stop()
	<-sub.Closed()
	return nil
}
func main() {
	_ = godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
