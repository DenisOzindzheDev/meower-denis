package db

import (
	"context"
	"meower-denis/internal/models"
)

type DB interface {
	Close()
	InsertMeow(ctx context.Context, meow models.Meow) error
	ListMeows(ctx context.Context, skip uint64, take uint64) ([]models.Meow, error)
}

var impl DB

func SetRepository(repository DB) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertMeow(ctx context.Context, meow models.Meow) error {
	return impl.InsertMeow(ctx, meow)
}
func ListMeows(ctx context.Context, skip uint64, take uint64) ([]models.Meow, error) {
	return impl.ListMeows(ctx, skip, take)
}
