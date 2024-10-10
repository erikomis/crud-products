package main

import (
	"goravel/controller"
	"goravel/db"
	"goravel/repository"
	"goravel/usecase"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	dbConnection, err := db.ConnectDB()

	if err != nil {
		panic(err)
	}

	defer dbConnection.Close()

	migrationDir := "./db/migrations"

	db.RunMigrations(dbConnection, migrationDir)
	ProductRepository := repository.NewProductRepository(dbConnection)

	ProductUseCase := usecase.NewProductUseCase(
		ProductRepository,
	)

	ProductController := controller.NewProductController(
		ProductUseCase,
	)

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/products", ProductController.GetProduct)
	server.POST("/product", ProductController.CreateProduct)
	server.GET("/product/:id", ProductController.GetProductById)
	server.DELETE("/product/:id", ProductController.DeleteProduct)
	server.PUT("/product/:id", ProductController.UpdateProduct)
	server.Run(":8888")
}
