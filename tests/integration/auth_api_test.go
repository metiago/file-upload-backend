package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/metiago/zbx1/common/request"
	"github.com/metiago/zbx1/repository"
)

var authPathURL = "login/auth"

func init() {
	mountBackEndURL()
}

func TestAuthUser(t *testing.T) {

	u := repository.User{
		Name:     "",
		Email:    "",
		Role:     &repository.Role{Name: ""},
		Username: "metiago",
		Password: "zero",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	status := request.PostHTTP(fmt.Sprintf("%s/%s", baseURL, authPathURL), "", data)

	if status != 200 {
		t.Errorf("Status was: %d", status)
	}
}
