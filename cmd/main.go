package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /{short_url}", useURL)
	router.HandleFunc("GET /api/urls/{short_url}", getURL)
	router.HandleFunc("POST /api/urls", createURL)
	router.HandleFunc("DELETE /api/urls/{short_url}", deleteURL)

	log.Println("Server up and listening on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}

var urls = make(map[string]URL)

type URL struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func useURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("short_url")

	entry, ok := urls[shortURL]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, entry.LongURL, http.StatusFound)
}

func getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("short_url")

	entry, ok := urls[shortURL]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

func createURL(w http.ResponseWriter, r *http.Request) {
	var url URL

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := urls[url.ShortURL]
	if ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	urls[url.ShortURL] = url

	w.WriteHeader(http.StatusCreated)
}

func deleteURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("short_url")

	_, ok := urls[shortURL]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(urls, shortURL)

	w.WriteHeader(http.StatusNoContent)
}
