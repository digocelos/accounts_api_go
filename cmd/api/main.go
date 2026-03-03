package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/digocelo/account-api/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	runMigrationFlag bool
)

func init() {
	flag.BoolVar(&runMigrationFlag, "migrate", false, "execute migration")
}

func main() {
	flag.Parse()
	dbURL := os.Getenv("APP_DB_URL")
	if dbURL == "" {
		panic("APP_DB_URL is required")
	}

	ctx := context.Background()
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if runMigrationFlag {
		log.Println("running migration...")
		err := postgres.RunMigrations(ctx, db.Pool)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migration completed")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	log.Println("Server started...")

	http.ListenAndServe(":3000", r)

}
