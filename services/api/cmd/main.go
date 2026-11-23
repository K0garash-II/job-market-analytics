package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/malayegor9-dotcom/job-market-analytics/pkg/db"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/handler"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/postgres"
	"log"
	"os"
)

func run() error {
	pool, err := db.NewPostgresPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()
	r := gin.Default()
	handler.New(postgres.Reader{Pool: pool}).RegisterRoutes(r)
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	return r.Run(":" + port)
}
func main() {
	_ = godotenv.Load()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
