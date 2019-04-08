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

	token := authorize()

	fileToUpload := "/home/tiago/Desktop/document.pdf"

	status := request.PostMultiPart(fmt.Sprintf("%s/%s", baseURL, filesPathURL), token, fileToUpload)

	if status != 200 {
		t.Errorf("Status was: %d", status)
	}
}
