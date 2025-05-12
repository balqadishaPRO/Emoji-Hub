package main

import (
	"log"
	"os"

	"github.com/balqadishaPRO/Emoji-Hub/internal/handler"
	"github.com/balqadishaPRO/Emoji-Hub/internal/middleware"
	"github.com/balqadishaPRO/Emoji-Hub/internal/repo"
	"github.com/balqadishaPRO/Emoji-Hub/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	repo, err := repo.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	svc := &service.EmojiService{Repo: repo}

	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		api.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://balqadishapro.github.io", "http://localhost:8080"},
			AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * 3600, // 12 hours
		}))
		api.Use(middleware.Session())
		handler.Register(api, svc)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(r.Run(":" + port))
}
