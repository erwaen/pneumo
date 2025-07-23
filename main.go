package main

import (
	"fmt"
	"github.com/erwaen/pneumo/store"
	"log"
	"net/http"
	"strings"
)

type Pneumonic struct {
	FullUrl   string
	PneumoUrl string
}

// creates a new object with the full url and the shortened 'pneumonic' version
func NewPneumonic(fullUrl string) *Pneumonic {
	return &Pneumonic{
		FullUrl:   fullUrl,
		PneumoUrl: strings.Split(fullUrl, "/")[1],
	}
}

func (p *Pneumonic) store() {
	fmt.Printf("storing...")
}

func handler(w http.ResponseWriter, r *http.Request) {
	pneumonic := NewPneumonic(r.URL.Path)
	fmt.Fprintf(w, "Hi this is the url %s and this is the pneumo version %s :)", pneumonic.FullUrl, pneumonic.PneumoUrl)

}

func main() {
	fmt.Println("checking on valkey...")

	store := store.CreateStore()
	defer store.Close()


	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
