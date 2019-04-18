package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/metiago/zbx1/common/request"
	"github.com/metiago/zbx1/repository"
)

var baseURL string

var authPathURL = "auth/login"

var token string

func mountBackEndURL() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	baseURL = "http://" + host + ":" + port
	token = authorize()
}

func authorize() string {
	u := repository.User{
		Username: "metiago",
		Password: "zero",
	}
	data, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}
	_, body := request.PostHTTPBody(fmt.Sprintf("%s/%s", baseURL, authPathURL), "", data)
	m := make(map[string]string)
	json.Unmarshal(body, &m)
	return m["token"]
}
