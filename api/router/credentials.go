package router

import (
	"xapiens/pkg/credential"

	"github.com/gin-gonic/gin"
)

func CredentialRoutes(service credential.Service, route *gin.Engine) {

	api := route.Group("/api")
	v1 := api.Group("/v1/credential")
	var d struct{}
	v1.POST("/store", func(c *gin.Context) {
		response, err := service.InsertCredential(c)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully stored",
			"data":    response,
		})
		return
	})

	v1.GET("/fetch", func(c *gin.Context) {
		response, err := service.FetchCredentials()
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully retrieved",
			"data":    response,
		})
	})

	v1.GET("/detail/:id", func(c *gin.Context) {
		response, err := service.DetailCredential(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully retrieved",
			"data":    response,
		})
		return
	})

	v1.PUT("/update", func(c *gin.Context) {
		response, err := service.UpdateCredential(c)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully updated",
			"data":    response,
		})
		return
	})

	v1.DELETE("/delete/:id", func(c *gin.Context) {
		err := service.DeleteCredential(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully deleted",
			"data":    d,
		})
		return
	})

	v1.POST("/login", func(c *gin.Context) {
		response, err := service.Login(c)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials successfully stored",
			"data":    response,
		})
		return
	})

	v1.GET("/validate-token", func(c *gin.Context) {
		err := service.IsValidToken(c)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "invalid credentials",
				"data":    d,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "credentials valid",
			"data":    d,
		})
		return
	})

}
