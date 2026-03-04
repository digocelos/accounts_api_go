package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/digocelo/account-api/internal/account"
	"github.com/digocelo/account-api/internal/httpapi"
	"github.com/digocelo/account-api/internal/repository/postgres"
)

var (
	runMigrationFlag bool
)

const (
	dbURL = "APP_DB_URL"
)

func init() {
	flag.BoolVar(&runMigrationFlag, "migrate", false, "execute migration")
}

func main() {
	flag.Parse()
	dbURL := os.Getenv(dbURL)
	if dbURL == "" {
		log.Fatalf("%s is required\n", dbURL)
	}

	logg := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

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
		log.Println("migration applyed")
		return
	}

	svc := account.NewService(postgres.NewAccountRepo(db.Pool))
	router := httpapi.NewRouter(logg, svc)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logg.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("server error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logg.Info("shutting down")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctxShutdown)

}
