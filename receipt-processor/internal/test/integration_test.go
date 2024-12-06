package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"receipt-processor/internal/controller"
	"receipt-processor/internal/model"
	"receipt-processor/internal/store"
	"receipt-processor/internal/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setup() {

	os.Setenv("ENV", "test")

	store.InitializeStores()
	utils.InitLogger()
}

func TestIntegration(t *testing.T) {

	setup()

	r := gin.Default()
	r.POST("/receipts/process", controller.ProcessReceipt)
	r.GET("/receipts/:id/points", controller.GetPoints)

	receipt := model.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []model.Item{
			{ShortDesc: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDesc: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}

	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal receipt: %v", err)
	}

	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receiptJSON))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	receiptID, ok := response["id"].(string)
	if !ok {
		t.Fatal("Expected receipt ID to be in the response")
	}

	t.Logf("Receipt ID: %s", receiptID)

	getReq, err := http.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	getRR := httptest.NewRecorder()
	r.ServeHTTP(getRR, getReq)

	assert.Equal(t, http.StatusOK, getRR.Code)

	var pointsResponse map[string]interface{}
	err = json.Unmarshal(getRR.Body.Bytes(), &pointsResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal points response: %v", err)
	}

	t.Logf("Points: %v", pointsResponse["points"])

	points, ok := pointsResponse["points"].(float64)
	if !ok {
		t.Fatal("Expected points to be in the response")
	}

	assert.Equal(t, 20, int(points), "Expected points to be 20")
}
