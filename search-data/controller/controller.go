package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/gorilla/mux"
	"github.com/opensearch-project/opensearch-go"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

var tenantID string
var IDToken string
var AccessToken string

// var OS_client_map map[string]*opensearch.Client
var employee = map[string]int{"Mark": 10, "Sandy": 20,
	"Rocky": 30, "Rajiv": 40, "Kate": 50}

var client *opensearch.Client

func ReadRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	vars := mux.Vars(r)
	tenantID = vars["tenantID"]
	IDToken = vars["IDToken"]
	AccessToken = vars["AccessToken"]

	fmt.Println(tenantID)
	fmt.Println(IDToken)
	fmt.Println(AccessToken)

	CheckCredentials(employee, IDToken)
}

func CheckCredentials(client_map map[string]int, IDToken string) {
	var count int = 1
	for key, element := range client_map {
		if key == IDToken {
			fmt.Println("Key:", key, "=>", "Element:", element)
			//call function "readMap" with IDToken and get the relevant os client
			break
		} else {
			//fmt.Println("Credentials not available")
			count = 2
		}
	}
	if count == 2 {
		fmt.Println("Credentials not available")
		//call function "find identity pool id" of tenant using tenantID
		//call function "getCredentialsForIdentityPool" using Ip_ID
		//call function "Find_OSInstance_URL" using tenantID
		//call function "Create_OSClient" using OS_instance url and credentials
		//call function "addClient" and add the new client to the global map using IDToken as the key
		//call function "readMap" with IDToken and get the relevant os client
	}

}

func Find_IP_ID(tenantID string) string {
	//return identity pool id
	var IP_ID string = "ip_id1"
	fmt.Println("Identity pool id is:")
	return IP_ID
}

/*
func Get_AWS_Credentials(IP_ID string){

}
*/
func getCredentialsForIdentityPool(identityPoolID string) (*cognitoidentity.Credentials, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	svc := cognitoidentity.New(sess)

	params := &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: aws.String(identityPoolID),
	}

	resp, err := svc.GetCredentialsForIdentity(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for identity pool %s: %v", identityPoolID, err)
	}

	return resp.Credentials, nil
}

func Find_OSInstance_URL(tenantID string) string {
	fmt.Println("OS instance URL is:")
	var test_url string = "https://search-odasara-test-domain-stosx4jruhkebwxwkvfsyjdln4.eu-west-2.es.amazonaws.com"
	return test_url
}

func Create_OSClient(url string, credentials *cognitoidentity.Credentials) *opensearch.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	cfg.Region = "eu-west-2"
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	fmt.Printf("cfg.Credentials: %v\n", cfg.Credentials)
	signer, err := requestsigner.NewSigner(cfg)
	if err != nil {
		fmt.Println(err.Error())
		//return
	}

	client, err = opensearch.NewClient(opensearch.Config{
		Addresses: []string{url},
		Signer:    signer,
	})
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		//return
	}
	fmt.Println(client.Info())
	return client
}

// func addClient(cli *opensearch.Client, IDToken string){
func addClient(cli int, IDToken string) {
	employee[IDToken] = cli

}

// func readMap(IDToken string)*opensearch.Client{
func readMap(IDToken string) int {
	cli := employee[IDToken]
	return cli
}
