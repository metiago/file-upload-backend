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

	"github.com/metiago/zbx1/common/request"

	"github.com/metiago/zbx1/repository"
)

func roleFindAll(w http.ResponseWriter, r *http.Request) {
	us, err := repository.FindAllRoles()
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

func roleFindOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	u, err := repository.FindRoleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			request.Handle404(w)
			return
		}
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

func roleAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	var s *repository.Role
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		request.Handle400(w, err)
		return
	}

	_, err = repository.AddRole(s)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func roleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	var s *repository.Role
	if err := json.Unmarshal(body, &s); err != nil {
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
	_, err = repository.UpdateRole(id, s)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func roleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	err = repository.DeleteRole(id)
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
	w.WriteHeader(http.StatusNoContent)
}
