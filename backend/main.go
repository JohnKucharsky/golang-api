package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	err := godotenv.Load()
	dbConn := os.Getenv("DB_ADDRESS")
	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		fmt.Println(dbConn)
		return sql.Open(
			"postgres",
			dbConn,
		)
	}

	if len(dbConn) == 0 {
		dbConn = fmt.Sprintf(
			"postgres://postgres:%s@db:5432/data?sslmode=disable",
			string(bin),
		)
	}

	return sql.Open(
		"postgres",
		dbConn,
	)
}

func blogHandler(w http.ResponseWriter, _ *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	rows, err := db.Query("SELECT title FROM blog")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var titles []string
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		titles = append(titles, title)
	}
	json.NewEncoder(w).Encode(titles)
}

func main() {
	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", blogHandler)
	log.Fatal(
		http.ListenAndServe(
			":8000", handlers.LoggingHandler(os.Stdout, r),
		),
	)
}

func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id SERIAL, title VARCHAR)"); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if _, err := db.Exec(
			"INSERT INTO blog (title) VALUES ($1);",
			fmt.Sprintf("Blog post #%d", i),
		); err != nil {
			return err
		}
	}
	return nil
}
