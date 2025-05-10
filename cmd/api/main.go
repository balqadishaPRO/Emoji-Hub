// cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// TODO: /emoji , /emoji/:id
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
