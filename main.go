package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/erwaen/pneumo/pneumify"
	"github.com/erwaen/pneumo/store"
)

type Server struct {
	store *store.StorageService
}

type PneumifyRequest struct {
	URL string `json:url`
}

type PneumifyResponse struct {
	ShortURL string `json:short_url`
}

func getFullURL(r *http.Request) string {
	// Get the scheme (http or https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Construct the full URL
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
}

func (s *Server) handlerPneumify(w http.ResponseWriter, r *http.Request) {

	var req PneumifyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request: could not decode JSON", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "Bad Request: 'url' field is required", http.StatusBadRequest)
		return
	}

	urlToPneumify := req.URL
	pneumo := pneumify.PneumifyURL(urlToPneumify)

	err = s.store.StorePneumo(urlToPneumify, pneumo)
	if err != nil {
		log.Printf("ERROR storing pneumo: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", pneumo)

	resp := PneumifyResponse{
		ShortURL: shortUrl,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handlerGetURL(w http.ResponseWriter, r *http.Request) {
	pneumo := r.PathValue("pneumo")
	url, err := s.store.RetrieveFromPneumo(pneumo)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        url = "http://" + url
    }

	http.Redirect(w, r, url, 302)
}

func main() {
	fmt.Println("checking on valkey...")

	store := store.CreateStore()
	defer store.Close()

	server := &Server{
		store: store,
	}

	http.HandleFunc("POST /pneumify/", server.handlerPneumify)
	http.HandleFunc("GET /{pneumo}", server.handlerGetURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
