// internal/middleware/session.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sid")
		if err != nil {
			sid = uuid.New().String()

			// Create cookie with SameSite=None
			cookie := &http.Cookie{
				Name:     "sid",
				Value:    sid,
				MaxAge:   3600 * 24 * 7, // 1 week
				Path:     "/",
				Domain:   "emoji-hub-6odk.onrender.com",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode, // ← Correct way to set SameSite
			}

			http.SetCookie(c.Writer, cookie)
		}

		c.Set("sid", sid)
		c.Next()
	}
}
