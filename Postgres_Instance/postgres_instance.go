package postgresinstance

import (
	"context"
	"fmt"
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

	// defer pool.Close()

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

// CreateDatabase creates a new database for the news application
func CreateDatabase() error {

	dbName := config.GetEnvironmentVariable("NEWS_DB_NAME")

	// Attempt to create the database
	_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if err != nil {
		log.Printf("Database %s already exists.\n", dbName)
		return fmt.Errorf("unable to create database: %v", err)
	}

	log.Printf("Database %s created successfully.\n", dbName)
	return nil
}
