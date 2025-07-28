package main

import (
	"encoding/json"
	"fmt"
	"github.com/erwaen/pneumo/pneumify"
	"github.com/erwaen/pneumo/store"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	store *store.StorageService
}

type PneumifyRequest struct {
	URL string `json:"url"`
}

type PneumifyResponse struct {
	ShortURL string `json:"short_url"`
}

func getDomainUrl(r *http.Request) string {
	// Get the scheme (http or https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Construct the full URL
	return fmt.Sprintf("%s://%s", scheme, r.Host)
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

	// Construct the short URL using the domain of getfullurl

	shortUrl := getDomainUrl(r) + "/" + pneumo	


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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("checking on valkey...")

	store := store.CreateStore()
	defer store.Close()

	server := &Server{
		store: store,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /pneumify/", server.handlerPneumify)
	mux.HandleFunc("GET /{pneumo}", server.handlerGetURL)

	handler := corsMiddleware(mux)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
