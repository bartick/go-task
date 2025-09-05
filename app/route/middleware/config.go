package middleware

import (
	"github.com/bartick/go-task/app/model"
	"github.com/gin-gonic/gin"
)

func Config(db model.DBTX) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
