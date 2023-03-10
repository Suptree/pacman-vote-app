package main

import (
	"log"
	"net/http"
)

func main() {
	port := "80"

	http.Handle("/", http.FileServer(http.Dir("./"))) // ①
	log.Printf("Server listening on http://localhost:%s/", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
