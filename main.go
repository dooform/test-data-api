package main

import (
	"fmt"

	"github.com/Dooform/test-data-api/config"
	"github.com/Dooform/test-data-api/database"
	"github.com/Dooform/test-data-api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()
	database.Migrate()

	r := gin.Default()

	// Add CORS middleware
	corsConfig := cors.DefaultConfig()
	origins := config.GetCORSOrigins()
	if origins == nil {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = origins
	}
	r.Use(cors.New(corsConfig))

	r.GET("/list", handlers.ListBoundaries)
	r.GET("/query", handlers.QueryBoundaries)
	r.GET("/search", handlers.SearchBoundaries)

	fmt.Println("Server is running on localhost:7242")
	r.Run("127.0.0.1:7242")
}
