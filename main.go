package main

import (
	"fmt"
	"log"
	"os"

	"xapiens/api/router"
	"xapiens/database"
	"xapiens/pkg/credential"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("database connection err ", err)
	}

	// migrations
	database.Migrations(db)

	// repo
	credentialRepo := credential.NewCredentialRepository(db)
	// service
	credentialService := credential.NewCredentialService(credentialRepo)
	fmt.Println(credentialService)

	route := gin.Default()
	router.CredentialRoutes(credentialService, route)
	router.HealthCheckRoutes(route)

	// NO ROUTE
	route.NoRoute(func(c *gin.Context) {
		var d struct{}
		c.JSON(404, gin.H{"message": "Page not found (SVC AUTH)", "success": false, "data": d})
	})

	route.Run(":" + os.Getenv("PORT"))
}
