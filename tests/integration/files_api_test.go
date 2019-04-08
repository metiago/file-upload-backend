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

func TestFileUpload(t *testing.T) {

	fileToUpload := "/home/tiago/Desktop/todo.txt"

	status := request.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesPathURL), token, fileToUpload)

	expected := 200
	if status != 200 {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}
