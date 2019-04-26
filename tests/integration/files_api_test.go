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

var filesPathURL = "api/v1/files"
var filesUploadPathURL = "api/v1/files/upload"

func init() {
	mountBackEndURL()
}

func TestFileUploadThenOK(t *testing.T) {

	fileToUpload := "int-tests.txt"

	status := helper.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesUploadPathURL), token, fileToUpload)

	expected := 201
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestFileUploadThenError(t *testing.T) {

	fileToUpload := "int-tests.txt"

	status := helper.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesUploadPathURL), "", fileToUpload)

	expected := 403
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestDeleteFile(t *testing.T) {

	file := getAnyFile()
	id := strconv.Itoa(file.ID)

	status := helper.DeleteHTTP(fmt.Sprintf("%s/%s/%s", baseURL, filesPathURL, id), token)

	expected := 204
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func getAnyFile() *repository.File {
	body, _ := helper.GetHTTP(fmt.Sprintf("%s/%s/%s", baseURL, filesPathURL, "metiago"), token)
	var files []*repository.File
	var file *repository.File
	if err := json.Unmarshal(body, &files); err != nil {
		log.Fatal(err)
	}
	for _, v := range files {
		file = v
	}

	return file
}

// TODO Test pagination func
