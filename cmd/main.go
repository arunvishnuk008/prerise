package main

import (
	"crimson-sunrise.site/pkg/db"
	"crimson-sunrise.site/pkg/initializers"
	"crimson-sunrise.site/pkg/middleware"
	"crimson-sunrise.site/pkg/routes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting Prerise..")
	var err error = nil
	initializers.LoadEnvVariables()
	db.DB, err = initializeDB()
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing DB: %v", err.Error())
		}
	}(db.DB)
	mux := mapRoutes()
	logger := slog.Default()
	// Adding middleware http
	mux = middleware.Apply(mux,
		middleware.PanicRecovery(logger),
	)
	// cors
	options := cors.Options{
		AllowedOrigins: []string{"http://localhost:4200","https://crimson-sunrise.site"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders: []string{"*"},
	}
	mux = cors.New(options).Handler(mux)
	server := &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: mux,
	}
	log.Println("Prerise has risen...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start Prerise due to error %s", err.Error())
	}
}

func initializeDB() (*sql.DB, error) {
	datasourceName := os.Getenv("DATABASE_NAME")
	if datasourceName == "" {
		log.Print("DATABASE_NAME is empty, using default embedded database!")
		datasourceName = "SunriseDB.db"
	}
	database, err := sql.Open("sqlite3", datasourceName)
	return database, err
}

func mapRoutes() http.Handler {
	mux := http.NewServeMux()

	// blog post apis
	mux.HandleFunc("GET /prerise/posts", routes.GetAll)
	mux.HandleFunc("GET /prerise/posts/{id}", routes.GetPostByID)
	mux.HandleFunc("POST /prerise/posts", middleware.AuthMiddleware(http.HandlerFunc(routes.NewPost)))
	mux.HandleFunc("PUT /prerise/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(routes.UpdatePostByID)))
	mux.HandleFunc("DELETE /prerise/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(routes.DeletePostByID)))

	// hot take apis
	mux.HandleFunc("GET /prerise/takes", routes.GetAllHotTakes)
	mux.HandleFunc("GET /prerise/takes/{id}", routes.GetHotTakeByID)
	mux.HandleFunc("POST /prerise/takes", middleware.AuthMiddleware(http.HandlerFunc(routes.NewHotTake)))
	mux.HandleFunc("DELETE /prerise/takes/{id}", middleware.AuthMiddleware(http.HandlerFunc(routes.DeleteHotTakeByID)))

	// auth APIs
	mux.HandleFunc("POST /prerise/authenticate", routes.Login)
	mux.HandleFunc("POST /prerise/revoke", middleware.AuthMiddleware(http.HandlerFunc(routes.Logout)))

	return mux
}