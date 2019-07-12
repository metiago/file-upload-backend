package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/metiago/zbx1/api"
	"github.com/metiago/zbx1/common/env"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" || port == "" {
		log.Println("You have to export HOST and PORT environment variables.")
		os.Exit(1)
	}
	service := host + ":" + port
	router := api.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/templates/static/", http.FileServer(http.Dir("./static/"))))
	
	migrate()

	log.Printf("Listen on %s", service)
	log.Fatal(http.ListenAndServe(service, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// TODO Change it to a migration tool after tests
func migrate() {

	verifyBytes, err := ioutil.ReadFile("ddl/ddl.sql")
	if err != nil {
		log.Fatal(err)
	}

	db := env.GetConnection()

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(string(verifyBytes))

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	tx.Commit()

}
