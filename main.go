package main

import (
	"fmt"
	"github.com/erwaen/pneumo/store"
	"log"
	"net/http"
)

func getFullURL(r *http.Request) string {
    // Get the scheme (http or https)
    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    
    // Construct the full URL
    return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
}


func handler(w http.ResponseWriter, r *http.Request) {
	url:= getFullURL(r)
	fmt.Println("Received request for URL:", url)
	afterSlash := r.URL.Path
	fmt.Println("afterslash: ", afterSlash[1:])



}

func main() {
	fmt.Println("checking on valkey...")

	store := store.CreateStore()
	defer store.Close()


	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
