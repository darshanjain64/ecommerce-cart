package main

import (
	"ecommerce-cart/middleware"
	"ecommerce-cart/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port==""{
		port="8000"
	}
	
	router:= gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.ProductRoutes(router)
	routes.CartRoutes(router)
	routes.AddressRoutes(router)

	router.Run(":"+port)
}
