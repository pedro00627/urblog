package main

import (
	"log"
	"net/http"
)

func main() {
	server, err := InitializeServer()
	if err != nil {
		log.Fatal(err)
	}

	// Iniciar servidor
	addr := ":8080"
	log.Printf("Iniciando servidor en %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
