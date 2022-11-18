package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/romik1505/auth/internal/store"
)

func NewPostgresConnenction() store.Storage {
	connString, ok := os.LookupEnv("PG_DSN")
	if !ok {
		log.Fatalln("PG_DSN not set")
	}
	log.Println("postgres connecion: ", connString)

	con, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatalln("database connection err: %w", err)
	}

	return store.Storage{
		DB: con,
	}
}
