package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /{short_url}", useUrl)

	router.HandleFunc("GET /api/urls/{short_url}", getUrl)
	router.HandleFunc("POST /api/urls", createUrl)
	router.HandleFunc("DELETE /api/urls/{short_url}", deleteUrl)

	log.Println("Server up and listening on port 8080")
	http.ListenAndServe(":8080", router)
}

var urls = make(map[string]Url)

type Url struct {
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
}

func useUrl(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("short_url")

	entry, ok := urls[shortUrl]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	var url Url

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := urls[url.ShortUrl]
	if ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	urls[url.ShortUrl] = url

	w.WriteHeader(http.StatusCreated)
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("short_url")

	_, ok := urls[shortUrl]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(urls, shortUrl)

	w.WriteHeader(http.StatusNoContent)
}
