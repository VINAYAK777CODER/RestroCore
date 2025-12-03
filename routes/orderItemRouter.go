package routes

import( "github.com/gin-gonic/gin"
	controller	"github.com/VINAYAK777CODER/RestroCore/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderitems",controller.GetOrderItems())
	incomingRoutes.GET("/orderitems/:orderitem_id",controller.GetOrderItem())

	incomingRoutes.GET("/orderItems-order/:order_id",controller.GetOrderItemByOrder())
	
	incomingRoutes.POST("/orderitems",controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderitems/:orderitem_id",controller.UpdateOrderItem()) 
}