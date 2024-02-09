package search

import (
	"context"
	"log"
	"meower-denis/internal/elastic"
	"meower-denis/internal/models"
	"meower-denis/internal/stream"
)

func OnMeowCreated(m stream.MeowCreatedMessage) {
	meow := models.Meow{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := elastic.InsertMeow(context.Background(), meow); err != nil {
		log.Println(err)
	}

}
