package main

import (
	"log"
	"net/http"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/iampassos/goshort/internal/domain"
	"github.com/iampassos/goshort/internal/urls"
)

func main() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY,
		short_url TEXT NOT NULL UNIQUE,
		long_url TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

		FOREIGN KEY (user_id) REFERENCES users (id)
			ON DELETE CASCADE
		)
		`)

	router := http.NewServeMux()

	router.HandleFunc("GET /{short_url}", useUrl)

	urls.RegisterRoutes(router)

	log.Println("Server up and listening on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}

func useUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("short_url")

	entry, ok := domain.Urls[shortUrl]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, entry.LongUrl, http.StatusFound)
}
