package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/opensearch-project/opensearch-go"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

var x_ID string          //id-token
var Authorization string //access-token
var x_tenant string      //tenant-id

var OS_client_map map[string]*opensearch.Client

//var employee = map[string]int{"Mark": 10, "Sandy": 20,
//	"Rocky": 30, "Rajiv": 40, "Kate": 50}

var client *opensearch.Client

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	w.WriteHeader(http.StatusOK)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	w.WriteHeader(http.StatusOK)
}

func ReadRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
		if k == "X_id" {
			ID := r.Header["X_id"]
			x_ID = ID[0]
			fmt.Printf("x_ID is %q\n", x_ID)

		} else if k == "Authorization" {
			auth := r.Header["Authorization"]
			Authorization = auth[0]

		} else if k == "X_tenant" {
			tenant := r.Header["X_tenant"]
			x_tenant = tenant[0]
		}
	}

	var c bool = CheckCredentials(OS_client_map, x_ID)
	fmt.Println(c)
	if !c {
		var cred *cognitoidentity.Credentials = cognito_identiy_GetId(x_ID)
		fmt.Println(cred)
		var cli *opensearch.Client = Create_OSClient("https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com", cred)
		fmt.Println(cli)
		//addClient(cli, x_ID)

	}
	//get the client with relevant key
	//copy the http request body to another request
	/*originalBody, _ := ioutil.ReadAll(r.Body)
	newReq, _ := http.NewRequest("POST", "https://example.com", bytes.NewBuffer(originalBody)) */
	//AddDocument(w http.ResponseWriter, r *http.Request)
	//select whether the request is a create or a search
	//call adddocument() or searchdocument() with needed parameters extracted from the endpoints to send new request to OS instance
	//forward the OS instance response

}

func CheckCredentials(client_map map[string]*opensearch.Client, x_ID string) bool {
	var count int = 1
	for key, element := range client_map {
		if key == x_ID {
			fmt.Println("Key:", key, "=>", "Element:", element)
			//call function "readMap" with IDToken and get the relevant os client
			count = 2
			break
		}
	}

	return count == 2
}

func Find_IP_ID(x_tenant string) string {
	//return identity pool id
	var IP_ID string = "ip_id1"
	fmt.Println("Identity pool id is:")
	return IP_ID
}

func cognito_identiy_GetId(token string) *cognitoidentity.Credentials {

	sess, _ := session.NewSession()
	svc := cognitoidentity.New(sess)
	logins := make(map[string]*string)
	logins["cognito-idp.eu-west-2.amazonaws.com/eu-west-2_ypy2SeovU"] = &token
	input := &cognitoidentity.GetIdInput{
		AccountId:      aws.String("210226302225"),
		IdentityPoolId: aws.String("eu-west-2:31716ca8-a2d0-492e-91a7-7c3a261af441"),
		Logins:         logins,
	}
	result, err := svc.GetId(input)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(*result.IdentityId)
	}
	getCredentialsForIdentityInput := &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: result.IdentityId,
		Logins:     logins,
	}
	getCredentialsForIdentityOutput, err := svc.GetCredentialsForIdentity(getCredentialsForIdentityInput)
	if err != nil {
		log.Println(err.Error())
		fmt.Println("credentials not got successfully!")
	} else {
		log.Println(*getCredentialsForIdentityOutput)
		fmt.Println("credentials got successfully!")
	}
	return getCredentialsForIdentityOutput.Credentials
}

func Find_OSInstance_URL(tenantID string) string {
	fmt.Println("OS instance URL is:")
	var test_url string = "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com"
	return test_url
}

func Create_OSClient(url string, cre *cognitoidentity.Credentials) *opensearch.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"),
		//config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET_KEY", "TOKEN")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*cre.AccessKeyId, *cre.SecretKey, *cre.SessionToken)),
	)

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	fmt.Printf("cfg.Credentials: %v\n", cfg.Credentials)
	signer, err := requestsigner.NewSigner(cfg)
	if err != nil {
		log.Println(err.Error())
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/"},
		Signer:    signer,
	})
	if err != nil {
		log.Printf("err: %s\n", err.Error())
	}
	fmt.Println("os client created!")
	//OS_client_map[x_ID] = client
	//fmt.Println(client.Info())
	return client
}

func addClient(cli *opensearch.Client, IDToken string) {

	OS_client_map[IDToken] = cli

	for key, val := range OS_client_map {
		fmt.Println("Header field %q, Value %q\n", key, &val)
	}

}

func readMap(IDToken string) *opensearch.Client {

	cli := OS_client_map[IDToken]
	return cli
}

type Document struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func AddDocument(w http.ResponseWriter, r *http.Request) {

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

type SearchRequest struct {
	Query string `json:"query"`
}

/*
func SearchDocument(w http.ResponseWriter, r *http.Request) {
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
*/

func SearchDocument(w http.ResponseWriter, r *http.Request) {
	//var searchRequest SearchRequest

	// Get the index and _id from the URL
	//url := r.URL.String()
	//target_index := url[1:strings.Index(url, "/_search/")]

}
