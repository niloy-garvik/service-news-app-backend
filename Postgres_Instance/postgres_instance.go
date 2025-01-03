package postgresinstance

import (
	"context"
	"log"
	"service-news-app-backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

// CreateRDSInstance initializes a singleton instance of a PostgreSQL connection
func CreatePostgresInstance() *pgxpool.Pool {
	if pool != nil {
		return pool // Return existing instance if already created
	}

	// Build the connection string from environment variables
	connStr := config.GetEnvironmentVariable("POSTGRES_URL")

	var err error
	pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	defer pool.Close()

	// Verify connectivity by pinging the database
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error verifying RDS instance: %v\n", err)
	}

	log.Println("Postgres database instance is reachable.")
	return pool
}

// GetRDSInstance returns the existing RDS database client instance
func GetPostgresInstance() *pgxpool.Pool {
	return pool
}
