package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odasaraik/search-data/router"
)

//var client *opensearch.Client

func main() {
	/*
		cfg, err := config.LoadDefaultConfig(context.TODO())
		cfg.Region = "eu-west-2"
		if err != nil {
			log.Fatalf("failed to load configuration, %v", err)
		}
		fmt.Printf("cfg.Credentials: %v\n", cfg.Credentials)
		signer, err := requestsigner.NewSigner(cfg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		client, err = opensearch.NewClient(opensearch.Config{
			Addresses: []string{"https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/"},
			Signer:    signer,
		})
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
		fmt.Println(client.Info())
		//get_roles(client)   */

	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":8088", r))
	fmt.Println("Listening at port 8088 ...")

	//http.HandleFunc("/", addDocument)
	//http.ListenAndServe(":8000", nil)

}
