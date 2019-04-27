package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"

	"github.com/metiago/zbx1/common/helper"

	"github.com/metiago/zbx1/repository"
)

const (
	authHeader    string = "Authorization"
	maxSizeInByte int64  = 16000000
)

// FIXME When upload a folder it crashes
func fileUpload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxSizeInByte)
	if err := r.ParseMultipartForm(maxSizeInByte); err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}

	if token := r.Context().Value("token"); token != nil {

		claims, err := helper.ExtractTokenClaims(token.(string), verifyKey)

		var u repository.User
		mapstructure.Decode(claims["uinf"], &u)

		var f = new(repository.File)
		f.Name = handler.Filename
		f.Ext = filepath.Ext(handler.Filename)
		f.Data = buf.Bytes()

		if validErrs := validateForUpload(f); len(validErrs) > 0 {
			err := map[string]interface{}{"validationError": validErrs}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = repository.FileUpload(u, f)
		if err != nil {

			if err == repository.ErrFileExists {
				log.Println(err)
				helper.Handle400(w, err)
				return
			}

			log.Println(err)
			helper.Handle500(w, err)
			return
		}

	}

	helper.HandleSuccessMessage(w, http.StatusCreated, "File has been uploaded successfully")
}

func validateForUpload(f *repository.File) url.Values {

	errs := url.Values{}

	if f.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	if len(f.Name) < 3 || len(f.Name) > 150 {
		errs.Add("name", "The name field must be between 3-150 chars!")
	}

	return errs
}

func fileDownload(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		helper.Handle500(w, err)
		return
	}

	f, err := repository.FileDownload(id)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.Handle404(w)
			return
		}
		helper.Handle500(w, err)
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

	query := r.URL.Query()
	username := query.Get("username")
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		helper.Handle500(w, err)
		return
	}

	result, err := repository.FindaAllFilesByUsername(username, page)
	if err != nil {
		helper.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		helper.Handle500(w, err)
		return
	}
}

func fileDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		log.Println(err)
		helper.Handle500(w, err)
		return
	}
	err = repository.DeleteFile(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			helper.Handle404(w)
			return
		}
		helper.Handle500(w, err)
		return
	}

	helper.HandleSuccessMessage(w, http.StatusNoContent, "File has been deleted successfully")
}
