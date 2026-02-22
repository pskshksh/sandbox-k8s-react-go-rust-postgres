package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	defer DB.Close()

	h := NewHandler(DB)

	router := mux.NewRouter()
	router.HandleFunc("/healthz", h.liveness).Methods("GET")
	router.HandleFunc("/readyz", h.readiness).Methods("GET")
	router.HandleFunc("/requests", h.getRequests).Methods("GET")
	router.HandleFunc("/requests", h.insertRequest).Methods("POST")

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
