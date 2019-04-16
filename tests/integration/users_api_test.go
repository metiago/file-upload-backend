package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/metiago/zbx1/common/request"
	"github.com/metiago/zbx1/repository"
)

var usersPathURL = "api/v1/users"

func init() {
	mountBackEndURL()
}

func TestAddUser(t *testing.T) {

	r := repository.Role{ID: 1, Name: "ADMIN"}
	u := &repository.User{
		Name:     "AAA",
		Email:    "ziggy@gmail.com",
		Username: "bbb",
		Password: "doggy",
		Role:     &r,
		Created:  time.Now()}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	status := request.PostHTTP(fmt.Sprintf("%s/%s", baseURL, usersPathURL), token, data)
	expected := 201
	if status != expected {
		t.Errorf("Status was: %d", status)
	}
}

func TestUpdateUser(t *testing.T) {
	anyUser := getAnyUser()
	anyRole := getAnyRole()
	r := repository.Role{ID: anyRole.ID, Name: ""}
	u := &repository.User{
		Name:     "Ziggy Update",
		Email:    "ziggy_update@gmail.com",
		Username: "ziggy_update",
		Password: "doggy_update",
		Role:     &r,
		Created:  time.Now()}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	id := strconv.Itoa(anyUser.ID)
	status := request.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token, data)
	expected := 200
	if status != expected {
		t.Errorf("Status was: %d", status)
	}
}

func TestFindOneUser(t *testing.T) {

	user := getAnyUser()
	id := strconv.Itoa(user.ID)

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token)

	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != expected {
		t.Errorf("Status was: %d", status)
	}

	if u == nil {
		t.Errorf("You should have at least one user registered.")
	}

}

func TestFindAllUsers(t *testing.T) {

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, usersPathURL), token)

	var users []*repository.User
	if err := json.Unmarshal(body, &users); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != expected {
		t.Errorf("Status was: %d", status)
	}

	if len(users) == 0 {
		t.Errorf("You should have at least one user registered.")
	}
}

func TestDeleteUser(t *testing.T) {

	user := getAnyUser()
	id := strconv.Itoa(user.ID)

	status := request.DeleteHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token)

	expected := 204
	if status != expected {
		t.Errorf("Status was: %d", status)
	}
}

func getAnyUser() *repository.User {
	body, _ := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, usersPathURL), token)
	var users []*repository.User
	var user *repository.User
	if err := json.Unmarshal(body, &users); err != nil {
		log.Fatal(err)
	}
	for _, v := range users {
		user = v
	}
	return user
}
