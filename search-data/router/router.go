package router

import (
	"github.com/gorilla/mux"
	"github.com/odasaraik/search-data/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/readRequest/{tenantID}/{IDToken}/{AccessToken}", controller.ReadRequest).Methods("GET")

	return router
}
