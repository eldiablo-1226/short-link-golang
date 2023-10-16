package main

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"short-link/config"
	"short-link/handlers"
	"short-link/repsoitory"
)

func main() {

	figure.NewFigure("SHORT LINK", "", true).Print()

	// Init configuration
	err := config.Init()
	if err != nil {
		log.Fatalf("Could't open con file: %v\n", err)
	}

	// Connect to repository
	err = repsoitory.Init(config.AppConfig.DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Close after closing
	defer repsoitory.Close()

	r := mux.NewRouter()
	handlers.SetShortLinkRoutes(r)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
