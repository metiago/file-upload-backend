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
	"github.com/metiago/zbx1/common/request"

	"github.com/metiago/zbx1/repository"
)

func userFindAll(w http.ResponseWriter, r *http.Request) {

	us, err := repository.FindAllUsers()
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(us); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
}

func userFindOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	u, err := repository.FindUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			request.Handle404(w)
			return
		}
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
}

func userAdd(w http.ResponseWriter, r *http.Request) {

	// READ JSON REQUEST BODY
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	// CHECK FOR ERROR
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	// CLOSE REQUEST BODY
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	// UNMARSHAL BODY INTO A STRUCTURE
	var s *repository.User
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		request.Handle400(w, err)
		return
	}

	vals := helper.Validate(s)
	if len(vals) > 0 {
		json.NewEncoder(w).Encode(vals)
		return
	}

	// ADD TO DATABASE
	_, err = repository.AddUser(s)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func userUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		return
	}
	var u *repository.User
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(err)
		request.Handle400(w, err)
		return
	}
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	_, err = repository.UpdateUser(id, u)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	err = repository.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			request.Handle404(w)
			return
		}
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
