package api

import (
	"net/http"
	"fmt"	
)

// Index handle index.html page
func index(w http.ResponseWriter, r *http.Request) {	
 	if r.Method == "GET" {		
		fmt.Fprintln(w, "API HEALTH OK")		
 	}
 }
