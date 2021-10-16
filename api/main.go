package main

import (
	"log"

	"xapiens/api/router"
	"xapiens/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("database connection err $s", err)
	}

	// migrations
	database.Migrations(db)

	// repo

	// service

	route := gin.Default()
	router.CredentialRoutes(route)
}
