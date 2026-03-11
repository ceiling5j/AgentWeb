package main

import (
	"log"
	"net/http"

	handler "online-agent-backend/api"
)

func main() {
	http.HandleFunc("/api/chat", handler.Handler)
	log.Println("backend listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
