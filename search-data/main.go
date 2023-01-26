package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/opensearch-project/opensearch-go"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

/*
	type Role_details struct {
		Reserved           bool     `json:"reserved"`
		Hidden             bool     `json:"hidden"`
		ClusterPermissions []string `json:"cluster_permissions"`
		IndexPermissions   []struct {
			IndexPatterns  []string `json:"index_patterns"`
			Fls            []string `json:"fls"`
			MaskedFields   []string `json:"masked_fields"`
			AllowedActions []string `json:"allowed_actions"`
		} `json:"index_permissions"`
		TenantPermissions []string `json:"tenant_permissions"`
		Static            bool     `json:"static"`
	}
*/
var client *opensearch.Client

type Document struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

/*
func get_roles(client *opensearch.Client) {
	body := strings.NewReader(`{}`)
	httpreq, err := http.NewRequest("GET", "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/_plugins/_security/api/roles/", body)
	httpreq.Header.Add("Content-Type", "application/json")
	resp, err := client.Perform(httpreq)
	if err != nil {
		log.Print(err)
		return
	}
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
	}
	v := map[string]Role_details{}
	unmarsh_err := json.Unmarshal(p, &v)
	if unmarsh_err != nil {
		log.Print(unmarsh_err.Error())
		return
	}
	log.Print(v["alerting_full_access"])
}*/

func main() {

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

	/*	httpreq, err := http.NewRequest("GET", "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/_plugins/_security/api/roles/", body)
		httpreq.Header.Add("Content-Type", "application/json")
		resp, err := client.Perform(httpreq)
		if err != nil {
			log.Print(err)
			return
		}*/

	client, err = opensearch.NewClient(opensearch.Config{
		Addresses: []string{"https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/"},
		Signer:    signer,
	})
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}
	fmt.Println(client.Info())
	//get_roles(client)

	http.HandleFunc("/", addDocument)
	http.ListenAndServe(":8000", nil)

	//http.HandleFunc("/", searchDocument)
	//http.ListenAndServe(":8000", nil)

	/*	settings := strings.NewReader(`{
		'settings': {
			'index': {
				'number_of_shards': 1,
				'number_of_replicas': 0
				}
			}
		}`)  */

	/*res := opensearchapi.IndicesCreateRequest{
	  	Index: "go-test-index1",
	  	Body:  settings,
	  }
	  fmt.Println("creating index", res.Pretty)  */
	/*
		index_name := "movies"

		// create an index
		if _, err := client.Indices.Create(index_name, client.Indices.Create.WithWaitForActiveShards("1")); err != nil {
			log.Fatal("indices.create", err)
		}
		fmt.Println("successfully created")

		// index a document
		document, _ := json.Marshal(map[string]interface{}{
			"title":    "Moneyball",
			"director": "Bennett Miller",
			"year":     "2011",
		})

		if _, err := client.Index(index_name, strings.NewReader(string(document)), client.Index.WithDocumentID(("1"))); err != nil {
			log.Fatal("index", err)
		}

		// wait for the document to index
		time.Sleep(1 * time.Second)
		fmt.Println("successfully indexed")

		// search for the document
		query, _ := json.Marshal(map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  "miller",
					"fields": []string{"title^2", "director"},
				},
			},
		})

		if resp, err := client.Search(client.Search.WithBody(strings.NewReader(string(query)))); err != nil {
			log.Fatal("index", err)
		} else {
			var r map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&r)
			hits := r["hits"].(map[string]interface{})["hits"]
			fmt.Println(hits)
		}    */

}

//http.HandleFunc("/", handleSearch)
//http.ListenAndServe(":8000", nil)
/*
func handleSearch(w http.ResponseWriter, r *http.Request) {
	// Get the full URL of the request
	url := r.URL.String()
	fmt.Println("URL: ", url)
	// Use string manipulation to extract the target index
	targetIndex := url[1:strings.Index(url, "/_search")]
	fmt.Println("Target index: ", targetIndex)
} */

func addDocument(w http.ResponseWriter, r *http.Request) {
	var document Document

	// Read the request body and unmarshal it into the document struct
	json.NewDecoder(r.Body).Decode(&document)

	// Marshal the document struct into json
	jsonValue, _ := json.Marshal(document)

	// Get the index and _id from the URL
	url := r.URL.String()
	index := url[1:strings.Index(url, "/_create/")]
	_id := url[strings.LastIndex(url, "/")+1:]

	// Create the endpoint
	endpoint := "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/" + index + "/_create/" + _id

	// Create a new request with the endpoint, json data and set the method to POST
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new http client
	//client := &http.Client{}
	resp, err := client.Perform(req)
	//resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("response Status:", resp.Status)
	// Print the response body
	fmt.Println("response Body:", resp.Body)
}

func searchDocument(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest

	// Read the request body and unmarshal it into the searchRequest struct
	json.NewDecoder(r.Body).Decode(&searchRequest)

	// Marshal the searchRequest struct into json
	//jsonValue, _ := json.Marshal(searchRequest)

	// Get the target index from the URL
	url := r.URL.String()
	targetIndex := url[1:strings.Index(url, "/_search")]

	// Create the endpoint
	endpoint := "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/" + targetIndex + "/_search"

	// Create a new request with the endpoint, json data and set the method to POST
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}

	/*req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	} */
	req.Header.Set("Content-Type", "application/json")

	// Create a new http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("response Status:", resp.Status)
	// Print the response body
	fmt.Println("response Body:", resp.Body)
}
