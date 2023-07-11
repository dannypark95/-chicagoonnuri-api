package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dannypark95/ChicagoOnnuri/config"
	"github.com/dannypark95/ChicagoOnnuri/routes"
	"github.com/gorilla/handlers"
)

func main() {
	// Load the environment variables
	config.LoadEnv()

	// Initialize the router
	router := routes.SetupRoutes()

	// Setup CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	originsOk := handlers.AllowedOrigins([]string{"https://www.chicagoonnuri.com"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
