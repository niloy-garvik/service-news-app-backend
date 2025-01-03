package main

import (
	"log"
	"net/http"

	PostgresInstance "service-news-app-backend/Postgres_Instance"
	envUtil "service-news-app-backend/config"
	"service-news-app-backend/routes" // Import the new routes package
)

func main() {
	// Load environment variables
	envUtil.LoadEnvironmentVariables()

	// Create database connection
	PostgresInstance.CreatePostgresInstance()

	// Setup routes
	r := routes.SetupRoutes()

	applicationPort := envUtil.GetEnvironmentVariable("applicationPort")

	log.Println("Server is starting...")

	err := http.ListenAndServe(":"+applicationPort, r)

	if err != nil {
		log.Fatalln("There's an error with the server,", err)
	} else {
		log.Println("Server is running at port:", applicationPort)
	}
}
