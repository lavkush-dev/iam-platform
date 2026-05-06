package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDB(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	// Connection pool settings (important for production)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("connected to postgres")

	return db
}
