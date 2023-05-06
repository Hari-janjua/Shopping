package Controllers

import (
	"Shopping/Models"
	services "Shopping/Services"
	"fmt"
	"log"
	"net/http"
	"time"

	guuid "github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	log.Println("IN: PlaceOrder()")
	defer log.Println("OUT: PlaceOrder()")

	var orderDetails Models.OrderDetails
	err := c.BindJSON(&orderDetails)
	if err != nil {
		log.Println("Error while Binding the data: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}
	fmt.Println("orderDetails: ", orderDetails)

	// Check that product count is less than 10
	err = services.CheckProductCountLessThanTen(orderDetails.Products)
	if err != nil {
		log.Println("Error while calling CheckProductCountLessThanTen Service: ", err)
		c.JSON(http.StatusInternalServerError, "COUNT_OF_PRODUCT_IS_MORE_THAN_10")
		return
	}

	// Extract the productIds from the orderDetails
	var productIds []string
	for _, product := range orderDetails.Products {
		productIds = append(productIds, product.ProductId)
	}

	// Get the product details by its productIds
	// We will fetch the productData from the database. So that user cannot send any unusual data according to his/her requirement
	productsDetail, err := services.GetProductsDetail(productIds)
	if err != nil {
		log.Println("Error while calling GetProductsDetail Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}
	// fmt.Println("products: ", productsDetail)

	var finalOrderData Models.Order
	premiumProductCount := 0
	totalPrice := 0
	for _, product := range productsDetail {
		// Get the count of premium products
		if product.Category.CategoryId == 1 {
			premiumProductCount += 1
		}

		// Get the ordered product Details from the orderedDetails
		orderedProduct, err := services.GetOrderedProductDetail(product, orderDetails)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "OUT_OF_STOCK")
		}

		// fmt.Println("orderedProduct: ", orderedProduct)

		productOrdered := Models.ProductModel{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			Price:       product.Price,
			Category: Models.Category{
				CategoryId:   product.Category.CategoryId,
				CategoryName: product.Category.CategoryName,
			},
			Count: orderedProduct.Count,
		}

		finalOrderData.Products = append(finalOrderData.Products, productOrdered)
		totalPrice += orderedProduct.Count * product.Price
	}

	// Assign the values to the Order struct
	id := guuid.New()
	// fmt.Printf("github.com/google/uuid:         %s\n", id.String())

	finalOrderData.OrderId = id.String()
	finalOrderData.TotalCost = totalPrice
	finalOrderData.GrandTotal = totalPrice
	finalOrderData.DiscountPercentage = 0

	// Discount of 10% will be given if products include 3 premium category products.
	if premiumProductCount >= 3 {
		finalOrderData.GrandTotal = totalPrice - (totalPrice * (1 / 10))
		finalOrderData.DiscountPercentage = 10
	}

	finalOrderData.CustomerId = orderDetails.CustomerId
	finalOrderData.Address = orderDetails.Address
	finalOrderData.OrderStatus = "Placed"
	// finalOrderData.DispatchDate =
	finalOrderData.CreatedBy = orderDetails.CustomerName
	finalOrderData.CreatedOn = time.Now()
	finalOrderData.ModifiedBy = orderDetails.CustomerName
	finalOrderData.ModifiedOn = time.Now()

	// fmt.Printf("%#v\n", finalOrderData)

	// Update the product catalog i.e update the count in stock
	err = services.UpdateProductCatalog(finalOrderData.Products, productsDetail)
	if err != nil {
		log.Println("Error while calling updateProductCatalog Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	err = services.InsertOrderDetails(finalOrderData)
	if err != nil {
		log.Println("Error while calling InsertOrderDetails Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, "Ordered placed successfully")

}

func GetOrderDetailsById(c *gin.Context) {
	log.Println("IN: PlaceOrder()")
	defer log.Println("OUT: PlaceOrder()")

	orderId := c.Param("id")

	orderDetail, err := services.GetOrderDetailsById(orderId)
	if err != nil {
		log.Println("Error while calling GetOrderDetailsById Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, orderDetail)
}

func GetAllOrdersDetailsOfSpecificUser(c *gin.Context) {
	log.Println("IN: GetAllOrdersDetailsOfSpecificUser()")
	defer log.Println("OUT: GetAllOrdersDetailsOfSpecificUser()")

	userId := c.Param("id")

	ordersDetail, err := services.GetOrderDetailsByUserId(userId)
	if err != nil {
		log.Println("Error while calling GetOrderDetailsById Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, ordersDetail)
}

func DispatchOrderById(c *gin.Context) {
	log.Println("IN: DispatchOrderById()")
	defer log.Println("OUT: DispatchOrderById()")

	orderId := c.Param("id")
	err := services.DispatchOrderByOrderId(orderId)
	if err != nil {
		log.Println("Error while calling DispatchOrderByOrderId Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, "Order Dispatched Successfully.")
}

func CompleteOrderById(c *gin.Context) {
	log.Println("IN: CompleteOrderById()")
	defer log.Println("OUT: CompleteOrderById()")

	orderId := c.Param("id")
	err := services.CompleteOrderByOrderId(orderId)
	if err != nil {
		log.Println("Error while calling CompleteOrderByOrderId Service: ", err)
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, "Order Completed.")
}
