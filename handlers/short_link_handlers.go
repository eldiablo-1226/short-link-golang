package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"short-link/config"
	"short-link/models"
	"short-link/repsoitory"
	"strings"
)

func redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	url, err := repsoitory.GetCodeByUrl(code)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func createShortLink(w http.ResponseWriter, r *http.Request) {

	var model models.InsertShortLink
	var code string
	var err error

	err = json.NewDecoder(r.Body).Decode(&model)

	if err != nil || model.Url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code, err = repsoitory.GetCodeByUrl(model.Url)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if code == "" {
		newUUID, _ := uuid.NewUUID()

		code = strings.ReplaceAll(newUUID.String(), "-", "")

		err = repsoitory.InsertShortLink(code, model.Url, model.Tag)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	result := models.Result{Url: fmt.Sprintf("%s/%s", config.AppConfig.Domain, code)}

	_ = json.NewEncoder(w).Encode(result)

	log.Printf("Short link created: %s\n", code)
}

func SetShortLinkRoutes(r *mux.Router) {
	r.HandleFunc("/{code}", redirectToOriginalURL).Methods("GET")
	r.HandleFunc("/shorter", Authorize(createShortLink)).Methods("POST")
}
