package routes

import (
	"receipt-processor/internal/controller"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	router := gin.Default()

	router.POST("/receipts/process", controller.ProcessReceipt)
	router.GET("/receipts/:id/points", controller.GetPoints)

	return router
}
