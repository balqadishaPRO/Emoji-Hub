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

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	// API routes
	api := r.Group("/api")
	{
		api.Use(middleware.Session())
		api.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:8080"},
			AllowMethods:     []string{"GET", "POST", "DELETE"},
			AllowCredentials: true,
		}))
		handler.Register(api, svc)
	}

	// Serve HTML pages
	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})
	r.GET("/catalog.html", func(c *gin.Context) {
		c.File("static/catalog.html")
	})
	r.GET("/favorites.html", func(c *gin.Context) {
		c.File("static/favorites.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(r.Run(":" + port))
}
