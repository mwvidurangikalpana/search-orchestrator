package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/gorilla/mux"
	"github.com/opensearch-project/opensearch-go"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

var x_ID string          //id-token
var Authorization string //access-token
var x_tenant string      //tenant-id

var OS_client_map map[string]*opensearch.Client

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

func AddDocument(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	vars := mux.Vars(r)

	GetHeaders(r.Header)

	//check whether credentials are availbale in the map
	var c bool = CheckCredentials(OS_client_map, x_ID)

	//if client is not there crate one
	if !c {
		SetAClient(x_ID)
	}
	//get the client with relevant key
	var read_cli *opensearch.Client = readMap(x_ID)

	index := vars["index"]
	_id := vars["_id"]

	// Create the endpoint
	endpoint := "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/" + index + "/_create/" + _id

	originalBody, _ := ioutil.ReadAll(r.Body)
	newReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(originalBody))

	if err != nil {
		panic(err)
	}

	newReq.Header.Set("Content-Type", "application/json")

	resp, err := read_cli.Perform(newReq)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	// Print the response status
	fmt.Println("response Status:", resp.Status)
	// Print the response body
	fmt.Println("response Body:", resp.Body)

	// Copy headers and status code from the new response to the original response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Copy the body from the new response to the original response
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func SearchDocument(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	vars := mux.Vars(r)

	GetHeaders(r.Header)

	//check whether credentials are availbale in the map
	var c bool = CheckCredentials(OS_client_map, x_ID)

	//if client is not there crate one
	if !c {
		SetAClient(x_ID)
	}
	//get the client with relevant key
	var read_cli *opensearch.Client = readMap(x_ID)

	target_index := vars["target-index"]

	endpoint := "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com/" + target_index + "/_search/"

	originalBody, _ := ioutil.ReadAll(r.Body)
	newReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(originalBody))

	if err != nil {
		panic(err)
	}

	newReq.Header.Set("Content-Type", "application/json")

	result, err := read_cli.Perform(newReq)
	if err != nil {
		panic(err)
	}

	defer result.Body.Close()
	// Print the response status
	fmt.Println("response Status:", result.Status)
	// Print the response body
	fmt.Println("response Body:", result.Body)

	// Copy headers and status code from the new response to the original response
	for k, v := range result.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(result.StatusCode)

	// Copy the body from the new response to the original response
	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetHeaders(body map[string][]string) {
	for k, v := range body {
		fmt.Printf("Header field %q, Value %q\n", k, v)
		if k == "X_id" {
			ID := body["X_id"]
			x_ID = ID[0]

		} else if k == "Authorization" {
			auth := body["Authorization"]
			Authorization = auth[0]

		} else if k == "X_tenant" {
			tenant := body["X_tenant"]
			x_tenant = tenant[0]
		}
	}
}

func SetAClient(x_ID string) {
	var cred *cognitoidentity.Credentials = cognito_identiy_GetId(x_ID)
	var cli *opensearch.Client = Create_OSClient("https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com", cred)
	addClient(cli, x_ID)
}

func CheckCredentials(client_map map[string]*opensearch.Client, x_ID string) bool {
	var count int = 1
	for key, element := range client_map {
		if key == x_ID {
			fmt.Println("Key:", key, "=>", "Element:", element)
			count = 2
			break
		}
	}

	return count == 2
}

// This function needs to be implemented with the service 'find identity pool id of the tenant using x-tenant header'
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

	} else {
		log.Println(*getCredentialsForIdentityOutput)

	}
	return getCredentialsForIdentityOutput.Credentials
}

// this function should be implemented using the service, 'find the opensearch instance url from opensearch provisioner service'
func Find_OSInstance_URL(tenantID string) string {
	fmt.Println("OS instance URL is:")
	var test_url string = "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com"
	return test_url
}

func Create_OSClient(url string, cre *cognitoidentity.Credentials) *opensearch.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*cre.AccessKeyId, *cre.SecretKey, *cre.SessionToken)),
	)

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

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

	return client
}

func addClient(cli *opensearch.Client, IDToken string) {

	OS_client_map = make(map[string]*opensearch.Client)
	OS_client_map[IDToken] = cli

}

func readMap(IDToken string) *opensearch.Client {

	cli := OS_client_map[IDToken]
	return cli
}
