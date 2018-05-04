package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pulpfree/gales-fuelsale-report/model"
)

// MiddleDB function
func MiddleDB(mongo model.DBHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mongo", mongo)
		c.Next()
	}
}
