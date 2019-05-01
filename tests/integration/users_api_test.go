package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/metiago/zbx1/common/helper"
	"github.com/metiago/zbx1/repository"
)

var usersPathURL = "api/v1/users"
var signUpURL = "signup"

func init() {
	mountBackEndURL()
}

func TestAddUser(t *testing.T) {

	u := &repository.User{
		Name:            "Bapi",
		Email:           "bapi@gmail.com",
		Username:        "bapi",
		Password:        "123XFS",
		ConfirmPassword: "123XFS",
		UpdatedPassword: ""}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	status := helper.PostHTTP(fmt.Sprintf("%s/%s", baseURL, signUpURL), "", data)
	expected := 201
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateUserPasswordThanOK(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy Update",
		Email:           "ziggy_update@gmail.com",
		Username:        anyUser.Username,
		Password:        "123XFS",
		ConfirmPassword: "123XFS",
		UpdatedPassword: "123XFS"}

	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	id := strconv.Itoa(anyUser.ID)
	up := string(id) + "/update-password"

	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, up), token, data)

	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateWrongConfirmPasswordThanFail(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy Update",
		Email:           "ziggy_update@gmail.com",
		Username:        anyUser.Username,
		Password:        "123XFS",
		ConfirmPassword: "Xsd",
		UpdatedPassword: "12345678"}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	id := strconv.Itoa(anyUser.ID)
	up := string(id) + "/update-password"

	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, up), token, data)

	expected := 400
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateWrongOldPasswordThanFail(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy Update",
		Email:           "ziggy_update@gmail.com",
		Username:        anyUser.Username,
		Password:        "kkk",
		ConfirmPassword: "kkk",
		UpdatedPassword: "12345678"}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	id := strconv.Itoa(anyUser.ID)
	up := string(id) + "/update-password"

	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, up), token, data)

	expected := 400
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateUser(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy",
		Email:           "ziggy@gmail.com",
		Username:        "ziggy",
		Password:        "12345678",
		ConfirmPassword: "12345678",
		UpdatedPassword: ""}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	id := strconv.Itoa(anyUser.ID)
	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token, data)
	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestAddExistingUserThanFail(t *testing.T) {

	u := &repository.User{
		Name:            "AAA",
		Email:           "ziggy@gmail.com",
		Username:        "metiago",
		Password:        "12345678",
		ConfirmPassword: "12345678",
		UpdatedPassword: ""}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	status := helper.PostHTTP(fmt.Sprintf("%s/%s", baseURL, signUpURL), "", data)

	expected := 400
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateExistingUserThanFail(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy Update",
		Email:           "ziggy_update@gmail.com",
		Username:        "metiago",
		Password:        "12345678",
		ConfirmPassword: "12345678",
		UpdatedPassword: ""}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	id := strconv.Itoa(anyUser.ID)
	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token, data)
	expected := 400
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateExistingUserThanOK(t *testing.T) {

	anyUser := getAnyUser()

	u := &repository.User{
		Name:            "Ziggy Update",
		Email:           "ziggy_update@gmail.com",
		Username:        anyUser.Username,
		Password:        "12345678",
		ConfirmPassword: "12345678",
		UpdatedPassword: ""}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}
	id := strconv.Itoa(anyUser.ID)
	status := helper.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token, data)
	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestFindOneUser(t *testing.T) {

	user := getAnyUser()
	id := strconv.Itoa(user.ID)

	body, status := helper.GetHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token)

	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}

	if u == nil {
		t.Errorf("You should have at least one user registered.")
	}

}

func TestFindAllUsers(t *testing.T) {

	body, status := helper.GetHTTP(fmt.Sprintf("%s/%s", baseURL, usersPathURL), token)

	var users []*repository.User
	if err := json.Unmarshal(body, &users); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}

	if len(users) == 0 {
		t.Errorf("You should have at least one user registered.")
	}
}

func TestDeleteUser(t *testing.T) {

	user := getAnyUser()
	id := strconv.Itoa(user.ID)

	status := helper.DeleteHTTP(fmt.Sprintf("%s/%s/%s", baseURL, usersPathURL, id), token)

	expected := 204
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func getAnyUser() *repository.User {
	body, _ := helper.GetHTTP(fmt.Sprintf("%s/%s", baseURL, usersPathURL), token)
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
