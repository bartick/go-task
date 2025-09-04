package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Config(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
