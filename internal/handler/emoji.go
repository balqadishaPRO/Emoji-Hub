package handler

import (
	"strconv"

	"github.com/balqadishaPRO/Emoji-Hub/internal/repo"
	"github.com/balqadishaPRO/Emoji-Hub/internal/service"
	"github.com/gin-gonic/gin"
)

func Register(r gin.IRouter, svc *service.EmojiService) {
	r.GET("/emoji", func(c *gin.Context) {
		p := repo.ListParams{
			Search:   c.Query("search"),
			Category: c.Query("category"),
			Group:    c.Query("group"),
			Sort:     c.DefaultQuery("sort", "name"),
			Limit:    2000,
			Offset:   0,
		}
		list, err := svc.List(c, p)
		if err != nil {
			c.JSON(500, errResp(err))
			return
		}
		c.JSON(200, list)
	})

	r.GET("/emoji/:id", func(c *gin.Context) {
		d, err := svc.Detail(c, c.Param("id"))
		if err != nil {
			c.JSON(404, errResp(err))
			return
		}
		c.JSON(200, d)
	})

	r.GET("/favorites", func(c *gin.Context) {
		sid := c.MustGet("sid").(string)
		out, err := svc.ListFav(c, sid)
		if err != nil {
			c.JSON(500, errResp(err))
			return
		}

		if len(out) == 0 {
			c.JSON(200, []interface{}{}) // Return empty array instead of null
			return
		}

		c.JSON(200, out)
	})

	r.POST("/favorites/:id", func(c *gin.Context) {
		err := svc.AddFav(c, c.MustGet("sid").(string), c.Param("id"))
		if err != nil {
			c.JSON(500, errResp(err))
			return
		}
		c.Status(204)
	})

	r.DELETE("/favorites/:id", func(c *gin.Context) {
		err := svc.DelFav(c, c.MustGet("sid").(string), c.Param("id"))
		if err != nil {
			c.JSON(500, errResp(err))
			return
		}
		c.Status(204)
	})

	r.OPTIONS("/emoji", handleOptions)
	r.OPTIONS("/favorites", handleOptions)
	r.OPTIONS("/favorites/:id", handleOptions)
}

func handleOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "https://balqadishapro.github.io")
	c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Status(204)
}

func atoi(s string) int       { n, _ := strconv.Atoi(s); return n }
func errResp(err error) gin.H { return gin.H{"error": err.Error()} }
