package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

	vals := helper.ValidateEmpty(u, "UpdatedPassword")
	if len(vals) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(vals)
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

	// TODO Return body with a custom message
	w.WriteHeader(http.StatusCreated)
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

	vals := helper.ValidateEmpty(u, "UpdatedPassword")
	if len(vals) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(vals)
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
	w.WriteHeader(http.StatusOK)
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

	u.ID = id
	_, err = repository.UpdateUserPassword(u)
	if err != nil {

		log.Println(err)

		switch err {
		case repository.ErrCheckPasswordEquality:
			helper.Handle400(w, err)
			return
		case repository.ErrMatchPassword:
			helper.Handle400(w, err)
			return
		default:
			helper.Handle500(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusNoContent)
}
