package api

import (
	"log"
	"meower-denis/internal/db"
	"meower-denis/internal/elastic"
	"meower-denis/internal/models"
	"meower-denis/internal/stream/nats"
	"meower-denis/pkg/util"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/segmentio/ksuid"
)

func CreateMeowHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}
	ctx := r.Context()

	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create meow")
		return
	}
	meow := models.Meow{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertMeow(ctx, meow); err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to insert meow")
		return
	}
	if err := nats.PublishMeowCreated(meow); err != nil {
		log.Println(err)
	}

	// Return new meow
	util.ResponseOK(w, response{ID: meow.ID})

}

func SearchMeowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")

	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		skip, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}

	meows, err := elastic.SearchMeows(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOK(w, []models.Meow{})
		return
	}
	util.ResponseOK(w, meows)

}

func ListMeowsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Read parameters
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}
	meows, err := db.ListMeows(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Cannot fetch meows")
	}
	util.ResponseOK(w, meows)
}
