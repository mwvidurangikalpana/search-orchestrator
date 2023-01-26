package router

import (
	"github.com/gorilla/mux"
	"github.com/odasaraik/search-data/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	//router.HandleFunc("/readRequest/{tenantID}/{IDToken}/{AccessToken}", controller.ReadRequest).Methods("GET")

	router.HandleFunc("/readRequest", controller.ReadRequest).Methods("GET")
	router.HandleFunc("/create/<index>/_create/<_id>", controller.AddDocument).Methods("POST")
	router.HandleFunc("/search/<target-index>/_search", controller.SearchDocument).Methods("POST")

	return router
}
