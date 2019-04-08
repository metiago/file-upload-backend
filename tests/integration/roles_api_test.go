package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/metiago/zbx1/common/request"
	"github.com/metiago/zbx1/repository"
)

const rolesPathURL = "api/v1/roles"

func init() {
	mountBackEndURL()
}

func TestAddRole(t *testing.T) {

	r := repository.Role{ID: 1, Name: "UNIT"}

	data, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}
	status := request.PostHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), token, data)

	expected := 201
	if status != 201 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestUpdateRole(t *testing.T) {
	any := getAnyRole()
	r := repository.Role{ID: any.ID, Name: "UPDATE"}
	data, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}
	status := request.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, strconv.Itoa(any.ID)), token, data)
	expected := 200
	if status != 200 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestFindOneRole(t *testing.T) {

	role := getAnyRole()
	id := strconv.Itoa(role.ID)

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, id), token)

	var u *repository.Role
	if err := json.Unmarshal(body, &u); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != 200 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}

	if u == nil {
		t.Errorf("You should have at least one Role registered.")
	}
}

func TestFindAllRoles(t *testing.T) {

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), token)

	var Roles []*repository.Role
	if err := json.Unmarshal(body, &Roles); err != nil {
		t.Error(err)
	}

	expected := 200
	if status != 200 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}

	if len(Roles) == 0 {
		t.Errorf("You should have at least one Role registered.")
	}
}

func TestDeleteRole(t *testing.T) {

	any := getAnyRole()
	anyID := strconv.Itoa(any.ID)

	status := request.DeleteHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, anyID), token)

	expected := 204
	if status != 204 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func getAnyRole() *repository.Role {
	body, _ := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), token)
	var roles []*repository.Role
	var role *repository.Role
	if err := json.Unmarshal(body, &roles); err != nil {
		log.Fatal(err)
	}
	for _, v := range roles {
		role = v
	}
	return role
}
