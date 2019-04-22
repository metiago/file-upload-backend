package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Jwt struct {
	AccessToken string `json:"token"`
}

// Response is a type to encapsulate http codes and messages
type Response struct {
	StatusCode int         `json:"status_code"`
	Message    interface{} `json:"message"`
}

func checkError(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}

func Handle200(w http.ResponseWriter, message string) {
	var response Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response.StatusCode = http.StatusForbidden
	response.Message = message
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle401(w http.ResponseWriter) {
	var response Response
	w.WriteHeader(http.StatusUnauthorized)
	response.StatusCode = http.StatusUnauthorized
	response.Message = "Unauthorized"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle403(w http.ResponseWriter) {
	var response Response
	w.WriteHeader(http.StatusForbidden)
	response.StatusCode = http.StatusForbidden
	response.Message = "Forbidden"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle406(w http.ResponseWriter) {
	var response Response
	w.WriteHeader(http.StatusNotAcceptable)
	response.StatusCode = http.StatusNotAcceptable
	response.Message = "Token expired"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle404(w http.ResponseWriter) {
	var response Response
	w.WriteHeader(http.StatusNotFound)
	response.StatusCode = http.StatusNotFound
	response.Message = "Data not found"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle400(w http.ResponseWriter, err error) {
	var response Response
	w.WriteHeader(http.StatusBadRequest)
	response.StatusCode = http.StatusBadRequest
	response.Message = fmt.Sprintf("%v", err)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Handle500(w http.ResponseWriter, err error) {
	var response Response
	w.WriteHeader(http.StatusInternalServerError)
	response.StatusCode = http.StatusInternalServerError
	response.Message = fmt.Sprintf("%v", err)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

// GetHTTP is a custom http get request
func GetHTTP(url string, token string) ([]byte, int) {
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode
}

// PostHTTP is a custom http post request
func PostHTTP(url string, token string, body []byte) int {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	return resp.StatusCode
}

func PostHTTPBody(url string, token string, body []byte) (int, []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	return resp.StatusCode, buf
}

// PostMultiPart is a custom http post request using multiparts
func PostMultiPart(url string, token string, file string) int {

	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", "name.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		log.Fatalln(err)
	}

	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	checkError(err)
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	return resp.StatusCode
}

// PutHTTP is a custom http put request
func PutHTTP(url string, token string, body []byte) int {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	return resp.StatusCode
}

// DeleteHTTP is a custom http put request
func DeleteHTTP(url string, token string) int {
	req, err := http.NewRequest("DELETE", url, nil)
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	return resp.StatusCode
}

// GetToken is used to get a access token in a back-end service
func GetToken(url string, username string, password string) string {
	credentials := map[string]interface{}{"username": username, "password": password}
	credentialsData, _ := json.Marshal(credentials)
	var jsonStr = []byte(string(credentialsData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var jwt Jwt
	json.Unmarshal(body, &jwt)
	return jwt.AccessToken
}
