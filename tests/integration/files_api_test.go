package tests

import (
	"fmt"
	"testing"

	"github.com/metiago/zbx1/common/request"
)

var filesPathURL = "api/v1/files/upload"

func init() {
	mountBackEndURL()
}

func TestFileUploadThenOK(t *testing.T) {

	fileToUpload := "/home/tiago/Desktop/todo.txt"

	status := request.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesPathURL), token, fileToUpload)

	expected := 200
	if status != 200 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestFileUploadThenError(t *testing.T) {

	fileToUpload := "/home/tiago/Desktop/todo.txt"

	status := request.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesPathURL), "", fileToUpload)

	expected := 403
	if status != 403 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}
