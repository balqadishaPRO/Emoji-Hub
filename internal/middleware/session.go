package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func Session() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if c.Request.Method == "OPTIONS" {
// 			c.Next() // Skip session checks for OPTIONS
// 			return
// 		}
// 		const name = "session_id"
// 		if sid, err := c.Cookie(name); err == nil {
// 			c.Set("sid", sid)
// 		} else {
// 			newID := uuid.New().String()
// 			c.SetCookie(name, newID, 3600*24*365, "/", "", false, true)
// 			c.Set("sid", newID)
// 		}
// 		c.Next()
// 	}
// }

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sid")
		if err != nil {
			sid = uuid.New().String()
			c.SetCookie(
				"sid",
				sid,
				3600*24*7,
				"/",
				"",
				true,
				true,
			)
		}
		c.Set("sid", sid)
		c.Next()
	}
}
