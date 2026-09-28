package handlers

import (
	"net/http"
	"strings"

	"shortner/internal/cache"
	"shortner/internal/database"
)

func GetLongURL(shortUrl string) string {
	output, err := cache.Get(shortUrl)
	if err == nil {
		return output
	}

	output, err = database.GetLongURL(shortUrl)
	if err != nil {
		return ""
	}

	cache.Set(shortUrl, output)

	return output
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	longURL := GetLongURL(shortURL)
	if longURL == "" {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
