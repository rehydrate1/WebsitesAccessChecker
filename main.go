package main

import (
	"log"
	"net/http"

	"github.com/rehydrate1/WebsitesAccessChecker/handlers"
)


func main() {
	fs := http.FileServer(http.Dir("./web/static"))

	http.Handle("/", fs)
	http.HandleFunc("/check", handlers.CheckHandler)
	
	port := ":8080"
	log.Printf("Starting server at port http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}