package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-processor/internal/model"
	"receipt-processor/internal/service"
	"receipt-processor/internal/utils"
	"receipt-processor/internal/validator"

	"github.com/gin-gonic/gin"
)

func ProcessReceipt(c *gin.Context) {
	fmt.Println("Processing receipt")
	data, err := utils.ParseFields(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var receipt model.Receipt
	updatedData, err := json.Marshal(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(updatedData, &receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := validator.ValidateReceipt(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := service.SaveReceipt(receipt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})

}

func GetPoints(c *gin.Context) {
	id := c.Param("id")
	points, err := service.GetPoints(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
}
