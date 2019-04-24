package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/metiago/zbx1/common/helper"

	"github.com/metiago/zbx1/repository"
)

func userFindAll(w http.ResponseWriter, r *http.Request) {

	us, err := repository.FindAllUsers()
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(us); err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}
}

func userFindOne(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])

	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	u, err := repository.FindUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			helper.Handle404(w)
			return
		}
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	// remove password hash for external requests
	u.Password = ""

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}
}

func userAdd(w http.ResponseWriter, r *http.Request) {

	// READ JSON REQUEST BODY
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	// CLOSE REQUEST BODY
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	// UNMARSHAL BODY INTO A STRUCTURE
	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(err)
		helper.Handle400(w, err)
		return
	}

	// VALIDATE REQUEST BODY FIELDS
	if validErrs := validateForInsert(u); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// ADD TO DATABASE
	_, err = repository.AddUser(u)
	if err != nil {

		if err == repository.ErrUsernameExists {
			log.Println(err)
			helper.Handle400(w, err)
			return
		}

		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	helper.HandleSuccessMessage(w, http.StatusCreated, "User has been created successfully")
}

func userUpdate(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Println(err)
		return
	}

	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(err)
		helper.Handle400(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return

	}

	if validErrs := validateForUpdate(u); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	u.ID = id
	_, err = repository.UpdateUser(u)
	if err != nil {

		if err == repository.ErrUsernameExists {
			log.Println(err)
			helper.Handle400(w, err)
			return
		}

		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	helper.HandleSuccessMessage(w, http.StatusOK, "User has been updated successfully")
}

func userUpdatePassword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Println(err)
		return
	}

	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(err)
		helper.Handle400(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	if validErrs := validateForUpdatePassword(u); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	u.ID = id
	_, err = repository.UpdateUserPassword(u)
	if err != nil {

		log.Println(err)

		switch err {
		case repository.ErrMatchPassword:
			helper.Handle400(w, err)
			return
		default:
			helper.Handle500(w, err)
			return
		}
	}

	helper.HandleSuccessMessage(w, http.StatusOK, "Your password has been updated successfully")
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}
	err = repository.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			helper.Handle404(w)
			return
		}
		helper.Handle500(w, err)
		return
	}

	helper.HandleSuccessMessage(w, http.StatusNoContent, "User has been deleted successfully")
}

func validateForInsert(u *repository.User) url.Values {

	errs := url.Values{}

	if u.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	if len(u.Name) < 3 || len(u.Name) > 120 {
		errs.Add("name", "The name field must be between 3-120 chars!")
	}

	if u.Email == "" {
		errs.Add("email", "The email field is required!")
	}

	if u.Username == "" {
		errs.Add("username", "The username field is required!")
	}

	if u.Password == "" {
		errs.Add("password", "The password field is required!")
	}

	if u.ConfirmPassword == "" {
		errs.Add("confirm_password", "The confirmation password field is required!")
	}

	if u.Password != u.ConfirmPassword {
		errs.Add("confirm_password", "Password and confirmation password not match!")
	}

	return errs
}

func validateForUpdate(u *repository.User) url.Values {

	errs := url.Values{}

	if u.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	if len(u.Name) < 3 || len(u.Name) > 120 {
		errs.Add("name", "The name field must be between 3-120 chars!")
	}

	if u.Email == "" {
		errs.Add("email", "The email field is required!")
	}

	if u.Username == "" {
		errs.Add("username", "The username field is required!")
	}

	return errs
}

func validateForUpdatePassword(u *repository.User) url.Values {

	errs := url.Values{}

	if u.Password == "" {
		errs.Add("password", "The password field is required!")
	}

	if u.ConfirmPassword == "" {
		errs.Add("confirm_password", "The confirmation password field is required!")
	}

	if u.UpdatedPassword == "" {
		errs.Add("updated_password", "The new password field is required!")
	}

	if u.ConfirmPassword != u.UpdatedPassword {
		errs.Add("confirm_password", "Password and confirmation password not match!")
	}

	return errs
}
