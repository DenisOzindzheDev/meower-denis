package api

import (
	"log"
	"meower-denis/internal/db"
	"meower-denis/internal/models"
	"meower-denis/internal/stream/nats"
	"meower-denis/pkg/util"
	"net/http"
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
