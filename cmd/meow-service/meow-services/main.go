package main

import (
	"fmt"
	"log"
	"meower-denis/internal/api"
	"meower-denis/internal/db"
	"meower-denis/internal/db/postgres"
	"meower-denis/internal/stream/nats"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	PostgresDB       string `envconfig :"POSTGRES_DB"`
	PostgresUser     string `envconfig :"POSTGRES_USER"`
	PostgresPassword string `envconfig :"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig :"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", api.CreateMeowHandler).Methods("POST").Queries("body", "{body}")
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		repo, err := postgres.NewPostgres(fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB))
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := nats.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		nats.SetEventStore(es)
		return nil
	})
	defer nats.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
