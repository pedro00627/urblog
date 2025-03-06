package main

import (
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	_ "github.com/pedro00627/urblog/docs"
)

func InitializeServer() (*http.ServeMux, error) {
	// Load dependencies
	deps, err := InitializeDependencies()
	if err != nil {
		log.Printf("Error al cargar dependencias: %v", err)
		return nil, err
	}

	// Configure routes
	mux := http.NewServeMux()
	ConfigureRoutes(mux, deps)

	// Serve swagger.yaml
	mux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../docs/swagger.yaml")
	})

	// Serve Swagger UI
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	mux.Handle("/docs", sh)

	return mux, nil
}
