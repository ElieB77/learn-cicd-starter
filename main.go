package main

import (
	"database/sql"
	"embed"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	DB *database.Queries
}

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Load .env file locally (only for local development)
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	// Hardcoded PORT to 8080
	port := "8080"

	apiCfg := apiConfig{}

	// Check if DATABASE_URL is set and connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
	}

	router := chi.NewRouter()

	// Set up CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Serve static files
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFiles.Open("static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Define API routes
	v1Router := chi.NewRouter()

	if apiCfg.DB != nil {
		v1Router.Post("/users", apiCfg.handlerUsersCreate)
		v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
		v1Router.Get("/notes", apiCfg.middlewareAuth(apiCfg.handlerNotesGet))
		v1Router.Post("/notes", apiCfg.middlewareAuth(apiCfg.handlerNotesCreate))
	}

	v1Router.Get("/healthz", handlerReadiness)

	// Mount API routes to the main router
	router.Mount("/v1", v1Router)

	// Create and start the server
	srv := &http.Server{
		Addr:              ":" + port, // Bind to the hardcoded port 8080
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
