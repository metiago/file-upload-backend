package tests

import (
	"encoding/json"
	"log"
	"os"

	"github.com/metiago/zbx1/common/request"
	"github.com/metiago/zbx1/repository"
)

var baseURL string

func mountBackEndURL() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	baseURL = "http://" + host + ":" + port
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
	_, body := request.PostHTTPBody(baseURL+"/login/auth", "", data)
	m := make(map[string]string)
	json.Unmarshal(body, &m)
	return m["token"]
}
