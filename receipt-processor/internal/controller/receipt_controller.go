package controller

import (
	"net/http"
	"receipt-processor/internal/model"
	"receipt-processor/internal/service"
	"receipt-processor/internal/utils"
	"receipt-processor/internal/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProcessReceipt(c *gin.Context) {
	utils.Logger.Info("Processing receipt request", zap.String("method", c.Request.Method), zap.String("url", c.Request.URL.Path))

	var requestPayload model.Receipt

	if err := validator.ValidateReceipt(requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		utils.Logger.Error("Invalid JSON format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	id, err := service.SaveReceipt(requestPayload)
	if err != nil {
		utils.Logger.Error("Failed to save receipt", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.Info("Receipt processed successfully", zap.String("receipt_id", id))
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetPoints(c *gin.Context) {

	id := c.Param("id")

	utils.Logger.Info("Fetching points", zap.String("receipt_id", id))

	points, err := service.GetPoints(id)
	if err != nil {
		utils.Logger.Error("Failed to retrieve points", zap.String("receipt_id", id), zap.Error(err))

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.Info("Points retrieved successfully", zap.String("receipt_id", id))

	c.JSON(http.StatusOK, gin.H{"points": points})
}
