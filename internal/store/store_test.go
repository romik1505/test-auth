package store

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var storage Storage

func TestMain(m *testing.M) {
	storage = NewPostgresConnenction()

	os.Exit(m.Run())
}

func NewPostgresConnenction() Storage {
	connString, ok := os.LookupEnv("PG_DSN")
	if !ok {
		log.Fatalln("PG_DSN not set")
	}
	log.Println("postgres connecion: ", connString)

	con, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatalln("database connection err: %w", err)
	}
	return Storage{
		DB: con,
	}
}

func mustTruncateAll() {
	storage.DB.Exec("DELETE FROM sessions") // nolint
	storage.DB.Exec("DELETE FROM users")    // nolint
}
