package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odasaraik/search-data/router"
)

func main() {

	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":8088", r))
	fmt.Println("Listening at port 8088 ...")

}
