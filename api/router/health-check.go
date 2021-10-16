package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckRoutes(route *gin.Engine) {

	route.GET("/", func(c *gin.Context) {
		var d struct{}
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Running SVC AUTH",
			"data":    d,
		})
	})
}
