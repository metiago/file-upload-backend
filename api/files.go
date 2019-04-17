package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"

	"github.com/metiago/zbx1/common/helper"
	"github.com/metiago/zbx1/common/request"

	"github.com/metiago/zbx1/repository"
)

const (
	authHeader    string = "Authorization"
	maxSizeInByte int64  = 16000000
)

func fileUpload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxSizeInByte)
	if err := r.ParseMultipartForm(maxSizeInByte); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	if token := r.Context().Value("token"); token != nil {

		claims, err := helper.ExtractTokenClaims(token.(string), verifyKey)

		var u repository.User
		mapstructure.Decode(claims["uinf"], &u)

		var f repository.File
		f.Name = handler.Filename
		f.Ext = filepath.Ext(handler.Filename)
		f.Data = buf.Bytes()
		err = repository.FileUpload(u, f)
		if err != nil {

			if err == repository.ErrFileExists {
				log.Println(err)
				request.Handle400(w, err)
				return
			}

			log.Println(err)
			request.Handle500(w, err)
		}

	}

	w.WriteHeader(http.StatusCreated)
}

func fileDownload(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		request.Handle500(w, err)
		return
	}

	f, err := repository.FileDownload(id)
	if err != nil {
		if err == sql.ErrNoRows {
			request.Handle404(w)
			return
		}
		request.Handle500(w, err)
		return
	}

	mime := http.DetectContentType(f.Data)

	contentLength := len(string(f.Data))

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	b := bytes.NewBuffer(f.Data)

	//stream the body to the client without fully loading it into memory
	io.Copy(w, b)
}

func fileFindAllByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	result, err := repository.FindaAllFilesByUsername(username)
	if err != nil {
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		request.Handle500(w, err)
		return
	}
}
