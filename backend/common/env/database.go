package env

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

var (
	settings string
)

func init() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	if host == "" || port == "" || username == "" || password == "" || database == "" {
		log.Fatal("You must export database environment variables.")
	}
	settings = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	log.Println("Connecting to database:", settings)
}

// GetConnection is responsible to get mysql connection
func GetConnection() *sql.DB {
	db, err := sql.Open("postgres", settings)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(80)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// InitRedisDB is responsible to get redis connection
func InitRedisDB() {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.3:6379",
		Password: "",
		DB:       0, // default DB
	})
	client.FlushDB()
}
