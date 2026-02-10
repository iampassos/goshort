package urls

import (
	"encoding/json"
	"net/http"

	"github.com/iampassos/goshort/internal/domain"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/urls/{short_url}", getUrl)
	router.HandleFunc("POST /api/urls", createUrl)
	router.HandleFunc("DELETE /api/urls/{short_url}", deleteUrl)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("short_url")

	entry, ok := domain.Urls[shortUrl]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	var url domain.Url

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := domain.Urls[url.ShortUrl]
	if ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	domain.Urls[url.ShortUrl] = url

	w.WriteHeader(http.StatusCreated)
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("short_url")

	_, ok := domain.Urls[shortUrl]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(domain.Urls, shortUrl)

	w.WriteHeader(http.StatusNoContent)
}
