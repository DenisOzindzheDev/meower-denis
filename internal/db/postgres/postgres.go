package postgres

import (
	"context"
	"database/sql"
	"log"
	"meower-denis/internal/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (r *Postgres) Close() {
	r.db.Close()
}

func (r *Postgres) InsertMeows(ctx context.Context, meow models.Meow) error {
	_, err := r.db.Exec("INSERT INTO meows (id, body, created_at) VALUES ($1, $2, $3)", meow.ID, meow.Body, meow.CreatedAt)
	return err
}

func (r *Postgres) ListMeows(ctx context.Context, skip uint64, take uint64) ([]models.Meow, error) {
	rows, err := r.db.Query("SELECT * FROM meows ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	meows := []models.Meow{}
	for rows.Next() {
		meow := models.Meow{}
		if err = rows.Scan(&meow.ID, &meow.Body, &meow.CreatedAt); err != nil {
			log.Fatal(err)
			if err = rows.Err(); err != nil {
				return nil, err
			}
		}
		meows = append(meows, meow)
	}
	return meows, nil
}
