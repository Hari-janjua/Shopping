package Routes

import (
	"Shopping/Controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.RouterGroup) {
	p := router.Group("/product")
	{
		p.POST("/create", Controllers.CreateNewProduct)
		p.GET("/getAllProduct", Controllers.GetAllProducts)
		p.GET("/:id", Controllers.GetProductById)
		p.DELETE("/delete/:id", Controllers.DeleteProduct)
	}
}
