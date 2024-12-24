package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdulazizax/yelp/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(config *config.Config) (*pgxpool.Pool, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.DBName,
	)

	// Create a context with a timeout for connecting to the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify the connection to the database
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	log.Printf("--------------------------- Connected to the database %s ---------------------------\n", config.Postgres.DBName)
	return pool, nil
}
