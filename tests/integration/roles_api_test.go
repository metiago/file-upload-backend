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
	status := request.PostHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), "", data)

	if status != 201 {
		t.Errorf("Status was: %d", status)
	}
}

func TestUpdateRole(t *testing.T) {
	any := getAnyRole()
	r := repository.Role{ID: any.ID, Name: "UPDATE"}
	data, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}
	status := request.PutHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, strconv.Itoa(any.ID)), "", data)
	if status != 200 {
		t.Errorf("Status was: %d", status)
	}
}

func TestFindOneRole(t *testing.T) {

	role := getAnyRole()
	id := strconv.Itoa(role.ID)

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, id), "")

	var u *repository.Role
	if err := json.Unmarshal(body, &u); err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Errorf("Status was: %d", status)
	}

	if u == nil {
		t.Errorf("You should have at least one Role registered.")
	}
}

func TestFindAllRoles(t *testing.T) {

	body, status := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), "")

	var Roles []*repository.Role
	if err := json.Unmarshal(body, &Roles); err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Errorf("Status was: %d", status)
	}

	if len(Roles) == 0 {
		t.Errorf("You should have at least one Role registered.")
	}
}

func TestDeleteRole(t *testing.T) {

	any := getAnyRole()
	anyID := strconv.Itoa(any.ID)

	status := request.DeleteHTTP(fmt.Sprintf("%s/%s/%s", baseURL, rolesPathURL, anyID), "")

	if status != 204 {
		t.Errorf("Status was: %d", status)
	}
}

func getAnyRole() *repository.Role {
	body, _ := request.GetHTTP(fmt.Sprintf("%s/%s", baseURL, rolesPathURL), "")
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
