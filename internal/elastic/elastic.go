package elastic

import (
	"context"
	"meower-denis/internal/models"
)

type Repository interface {
	Close()
	InsertMeow(ctx context.Context, meow models.Meow) error
	SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]models.Meow, error)
}

var impl Repository

func SetRepository(repositroy Repository) {
	impl = repositroy
}
func Close() {
	impl.Close()
}
func InsertMeow(ctx context.Context, meow models.Meow) error {
	return impl.InsertMeow(ctx, meow)
}
func SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]models.Meow, error) {
	return impl.SearchMeows(ctx, query, skip, take)
}
