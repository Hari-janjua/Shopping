package Routes

import (
	"Shopping/Controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoute(router *gin.RouterGroup) {
	o := router.Group("/order")
	{
		// API's for User
		o.POST("/placeOrder", Controllers.PlaceOrder)
		o.GET("/getOrderDetailsById/:id", Controllers.GetOrderDetailsById)
		o.GET("/getAllOrdersDetailsOfSpecificUser/:id", Controllers.GetAllOrdersDetailsOfSpecificUser)

		// API's for Admin
		o.PATCH("/dispatchOrder/:id", Controllers.DispatchOrderById)
		o.PATCH("/completeOrder/:id", Controllers.CompleteOrderById)

	}
}
