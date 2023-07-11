package routes

import (
	"log"
	"net/http"

	"github.com/dannypark95/ChicagoOnnuri/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	log.Println("Setting up routes...")
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/jubo", ShowJubo).Methods("GET") 
	router.HandleFunc("/login", Login).Methods("POST")

	// Protected routes
	router.Handle("/listPDF", middleware.AuthMiddleware(http.HandlerFunc(ListPDFs))).Methods("GET")
	router.Handle("/setLiveJubo", middleware.AuthMiddleware(http.HandlerFunc(SetLiveJubo))).Methods("POST")
	router.Handle("/pdf", middleware.AuthMiddleware(http.HandlerFunc(DeletePDF))).Methods("DELETE")
	router.Handle("/pdf", middleware.AuthMiddleware(http.HandlerFunc(UploadPDF))).Methods("POST")

	return router
}
