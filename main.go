package main

import (
	"context"
	"log"
	"net/http"

	PostgresInstance "service-news-app-backend/Postgres_Instance"
	schemas "service-news-app-backend/Schemas"
	envUtil "service-news-app-backend/config"
	"service-news-app-backend/routes" // Import the new routes package
)

func main() {
	// Load environment variables
	envUtil.LoadEnvironmentVariables()

	// Create database connection
	PostgresInstance.CreatePostgresInstance()
	PostgresInstance.CreateDatabase()
	schemas.CreateArticlesTable(context.Background(), PostgresInstance.GetPostgresInstance())

	// Setup routes
	r := routes.SetupRoutes()

	applicationPort := envUtil.GetEnvironmentVariable("applicationPort")

	log.Println("Server is running at port:", applicationPort)
	http.ListenAndServe(":"+applicationPort, r)

}
