package main

import (
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func InitializeServer() (*http.ServeMux, error) {
	// Cargar dependencias
	deps, err := InitializeDependencies()
	if err != nil {
		log.Printf("Error al cargar dependencias: %v", err)
		return nil, err
	}

	// Configurar rutas
	mux := http.NewServeMux()
	ConfigureRoutes(mux, deps)

	// Configurar Swagger
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../docs/swagger.yaml")
	})

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger"}
	sh := middleware.SwaggerUI(opts, nil)
	mux.Handle("/docs", sh)

	return mux, nil
}
