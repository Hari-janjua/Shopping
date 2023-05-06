package Controllers

import (
	"Shopping/Models"
	services "Shopping/Services"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateNewProduct(c *gin.Context) {
	log.Println("IN: CreateNewProduct")
	log.Println("OUT: CreateNewProduct")
	var product Models.Product

	err := c.BindJSON(&product)
	if err != nil {
		log.Println("Error while Binding the data: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	err = services.CreateNewProduct(product)
	if err != nil {
		log.Println("Error while calling CreateNewProduct Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, "Data Saved Successully")

}

func GetAllProducts(c *gin.Context) {
	log.Println("IN: GetAllProducts")
	defer log.Println("OUT: GetAllProducts")

	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))
	pageNumber -= 1

	var products []Models.Product
	products, err := services.GetAllProducts(pageNumber)
	if err != nil {
		log.Println("Error while calling service GetAllProducts(): ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	log.Println("IN: GetProductById")
	defer log.Println("OUT: GetProductById")

	productId := c.Param("id")

	product, err := services.GetProductById(productId)
	if err != nil {
		log.Println("Error while calling service GetProductById(): ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	fmt.Println("product: ", product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	log.Println("IN: GetProductById")
	defer log.Println("OUT: GetProductById")

	productId := c.Param("id")

	err := services.DeleteProductById(productId)
	if err != nil {
		log.Println("Error while calling service DeleteProductById(): ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, "Data removed successfully")

}
