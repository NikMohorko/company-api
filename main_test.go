package main_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var Token string
var Username string
var Port string
var CompanyId string

type authResponse struct {
	Token    string
	Validity int
}

func TestCreateUser(t *testing.T) {

	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		t.Fail()
		t.Logf("Failed to load environment variables: " + err.Error())
		return
	}

	Port = os.Getenv("HOST_PORT")

	httpClient := http.Client{}

	// Use a random username to prevent unique constraint violation when running tests multiple times
	maxVal := big.NewInt(100000000)
	randomAdd, _ := rand.Int(rand.Reader, maxVal)
	Username = "testing" + fmt.Sprint(randomAdd)

	body := map[string]string{
		"username": Username,
		"password": "testing",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "http://localhost:"+Port+"/user/create", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("User creation failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf("User creation failed: " + err.Error())
	}

	if resp.StatusCode != 201 {
		t.Fail()
		t.Logf("Expected status code 201, got: " + fmt.Sprint((resp.StatusCode)))
	}

}

func TestAuthenticateUser(t *testing.T) {

	httpClient := http.Client{}

	body := map[string]string{
		"username": Username,
		"password": "testing",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "http://localhost:"+Port+"/user/authenticate", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("User creation failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf("User authentication failed: " + err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fail()
		t.Logf("Expected status code 200, got: " + fmt.Sprint((resp.StatusCode)))
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	response := authResponse{}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	if response.Token == "" {
		t.Fail()
		t.Logf("Missing token in response")
	}

	Token = response.Token

}

func TestCreateCompany(t *testing.T) {

	httpClient := http.Client{}

	body := map[string]interface{}{
		"name":           "testCompany",
		"description":    "testing",
		"employee_count": 10,
		"registered":     true,
		"type":           "Cooperative",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "http://localhost:"+Port+"/company/create", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("Company creation failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf(err.Error())
	}

	if resp.StatusCode != 201 {
		t.Fail()
		t.Logf("Expected status code 201, got: " + fmt.Sprint((resp.StatusCode)))
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	var response map[string]interface{}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	companyUUID, ok := response["Id"]

	if !ok {
		t.Fail()
		t.Logf("Company ID not in response")
	}

	CompanyId = companyUUID.(string)

}

func TestUpdateCompany(t *testing.T) {

	httpClient := http.Client{}

	body := map[string]interface{}{
		"name":           "testCompany",
		"description":    "newTestingDescription",
		"employee_count": 10,
		"registered":     true,
		"type":           "Cooperative",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PATCH", "http://localhost:"+Port+"/company/update?id="+CompanyId, bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("Company update failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf(err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fail()
		t.Logf("Expected status code 200, got: " + fmt.Sprint((resp.StatusCode)))
	}

}

func TestGetCompany(t *testing.T) {

	httpClient := http.Client{}

	body := map[string]interface{}{}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("GET", "http://localhost:"+Port+"/company/get?name=testCompany", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("Company get failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf(err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fail()
		t.Logf("Expected status code 200, got: " + fmt.Sprint((resp.StatusCode)))
		return
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	var response map[string]interface{}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fail()
		t.Logf("Malformed response data: " + err.Error())
	}

	companyName, ok := response["description"]

	if !ok || companyName != "newTestingDescription" {
		t.Fail()
		t.Logf("Expected company description 'newTestingDescription'")
	}

}

func TestDeleteCompany(t *testing.T) {

	httpClient := http.Client{}

	body := map[string]interface{}{}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("DELETE", "http://localhost:"+Port+"/company/delete?name=testCompany", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fail()
		t.Logf("Company deletion failed: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := httpClient.Do(req)

	if err != nil {
		t.Fail()
		t.Logf(err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fail()
		t.Logf("Expected status code 200, got: " + fmt.Sprint((resp.StatusCode)))
	}

}
