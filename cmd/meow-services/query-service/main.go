package main

import (
	"fmt"
	"log"
	"meower-denis/internal/api"
	"meower-denis/internal/db"
	"meower-denis/internal/db/postgres"
	"meower-denis/internal/elastic"
	"meower-denis/internal/elastic/search"
	"meower-denis/internal/stream/nats"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", api.ListMeowsHandler).
		Methods("GET")
	router.HandleFunc("/search", api.SearchMeowsHandler).
		Methods("GET")
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
		repo, err := postgres.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.NewElastic(fmt.Sprintf("https://%s", cfg.ElasticsearchAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		elastic.SetRepository(es)
		return nil
	})
	defer elastic.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := nats.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		err = es.OnMeowCreated(search.OnMeowCreated)
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
