package search

import (
	"context"
	"encoding/json"
	"log"
	"meower-denis/internal/models"

	"github.com/olivere/elastic"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &ElasticRepository{client: client}, nil
}
func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertMeow(ctx context.Context, meow models.Meow) error {
	_, err := r.client.Index().
		Index("meows").
		Type("meow").
		Id(meow.ID).
		BodyJson(meow).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticRepository) SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]models.Meow, error) {
	result, err := r.client.Search().
		Index("meows").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	meows := []models.Meow{}
	for _, hit := range result.Hits.Hits {
		var meow models.Meow
		if err = json.Unmarshal(*hit.Source, &meow); err != nil {
			log.Println(err)
		}
		meows = append(meows, meow)
	}
	return meows, nil
}
