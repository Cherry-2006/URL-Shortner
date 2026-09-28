package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"shortner/internal/cache"
	"shortner/internal/database"
	"shortner/internal/shortener"
)

type ShortenRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

var (
	ErrInvalidAlias = errors.New("alias must be 1-32 characters using letters, numbers, underscore or hyphen")
	ErrAliasTaken   = errors.New("alias is already taken")
)

var aliasPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,32}$`)

var reservedAliases = map[string]bool{
	"shorten":     true,
	"app.js":      true,
	"style.css":   true,
	"favicon.ico": true,
}

func validAlias(alias string) bool {
	if !aliasPattern.MatchString(alias) {
		return false
	}
	if reservedAliases[strings.ToLower(alias)] {
		return false
	}
	return true
}

func aliasTaken(alias string) bool {
	if _, err := cache.Get(alias); err == nil {
		return true
	}
	if _, err := database.GetLongURL(alias); err == nil {
		return true
	}
	return false
}

func isUniqueError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE")
}

func GenAndInsertLong(longUrl string, alias string) (string, error) {
	alias = strings.TrimSpace(alias)
	if alias != "" {
		if !validAlias(alias) {
			return "", ErrInvalidAlias
		}
		if aliasTaken(alias) {
			return "", ErrAliasTaken
		}
		if err := database.InsertURL(longUrl, alias); err != nil {
			if isUniqueError(err) {
				return "", ErrAliasTaken
			}
			return "", err
		}
		cache.Set(alias, longUrl)
		return alias, nil
	}
	cnt := database.GetNextCounter()
	short := shortener.GenerateShort(cnt)
	if err := database.InsertURL(longUrl, short); err != nil {
		return "", err
	}
	cache.Set(short, longUrl)
	return short, nil
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortURL, err := GenAndInsertLong(req.URL, req.Alias)
	if err != nil {
		if errors.Is(err, ErrAliasTaken) {
			http.Error(w, ErrAliasTaken.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, ErrInvalidAlias) {
			http.Error(w, ErrInvalidAlias.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(ShortenResponse{
		ShortURL: shortURL,
	})
}
