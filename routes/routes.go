package routes

import (
	"github.com/gorilla/mux"
	"backend/controllers"
	"backend/middlewares"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsersHandler).Methods("GET")
	// Protected routes (require JWT authentication)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.JWTAuthMiddleware)
	// File upload route 
	protected.HandleFunc("/upload", controllers.UploadFile).Methods("POST")
	// File retrieval route (for retrieving file metadata)
	protected.HandleFunc("/files", controllers.FileMetadataHandler).Methods("GET")
	// File sharing route (for generating a public shareable link)
	protected.HandleFunc("/share/{file_id}", controllers.ShareFileHandler).Methods("GET")
	
	return router
}
