package main

import (
	"fmt"
	"github.com/erwaen/pneumo/minivalkey"
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

	// check if can connect to valkey
	valkeyClient, err := minivalkey.CreateClient("localhost:6379")

	if err != nil {
		log.Fatalf("Could not connect to valkey: %v", err)
	} else {
		fmt.Println("Connected to valkey successfully!")
	}
	defer valkeyClient.Close()

	// check health
	resp, err := valkeyClient.SendRespCommand("PING")
	if err != nil {
		log.Fatalf("Error sending PING command: %v", err)
	}
	fmt.Println("response", resp)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
