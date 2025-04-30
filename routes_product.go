package main

import (
	"github.com/gin-gonic/gin"
)

func registerProductRoutes(r *gin.Engine) {
	productGroup := r.Group("/products")
	{
		productGroup.GET("/", getAllProducts)
		productGroup.POST("/", createProduct)
	}
}

func getAllProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get all products"})
}

func createProduct(c *gin.Context) {
	c.JSON(201, gin.H{"message": "Create product"})
}