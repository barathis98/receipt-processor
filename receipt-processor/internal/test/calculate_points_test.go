package test

import (
	"receipt-processor/internal/model"
	"receipt-processor/internal/service"
	"testing"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name    string
		receipt model.Receipt
		want    int
		wantErr bool
	}{
		{
			name: "Base case 1",
			receipt: model.Receipt{
				Retailer: "Target",
				Total:    35.35,
				Items: []model.Item{
					{ShortDesc: "Mountain Dew 12PK", Price: 6.49},
					{ShortDesc: "Emils Cheese Pizza", Price: 12.25},
					{ShortDesc: "Knorr Creamy Chicken", Price: 1.26},
					{ShortDesc: "Doritos Nacho Cheese", Price: 3.35},
					{ShortDesc: "Klarbrunn 12-PK 12 FL OZ", Price: 12.00},
				},
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
			},
			want:    28,
			wantErr: false,
		},
		{
			name: "Base case 2",
			receipt: model.Receipt{
				Retailer: "M&M Corner Market",
				Total:    9.00,
				Items: []model.Item{
					{ShortDesc: "Gatorade", Price: 2.25},
					{ShortDesc: "Gatorade", Price: 2.25},
					{ShortDesc: "Gatorade", Price: 2.25},
					{ShortDesc: "Gatorade", Price: 2.25},
				},
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
			},
			want:    109,
			wantErr: false,
		},
		{
			name: "Valid receipt with integer total and divisible by 0.25",
			receipt: model.Receipt{
				Retailer:     "Store1",
				Total:        10.00,
				Items:        []model.Item{{ShortDesc: "Item1", Price: 5.0}, {ShortDesc: "Item2", Price: 3.0}},
				PurchaseDate: "2024-12-05",
				PurchaseTime: "14:30",
			},
			want:    102,
			wantErr: false,
		},
		{
			name: "Invalid total amount (negative)",
			receipt: model.Receipt{
				Retailer:     "Store2",
				Total:        -5.00,
				Items:        []model.Item{{ShortDesc: "Item1", Price: 5.0}},
				PurchaseDate: "2024-12-05",
				PurchaseTime: "15:00",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Valid receipt with odd day of the month",
			receipt: model.Receipt{
				Retailer:     "Store3",
				Total:        8.00,
				Items:        []model.Item{{ShortDesc: "Item1", Price: 2.0}, {ShortDesc: "Item2", Price: 3.0}},
				PurchaseDate: "2024-12-05",
				PurchaseTime: "14:15",
			},
			want:    102,
			wantErr: false,
		},
		{
			name: "Valid receipt with description length multiple of 3",
			receipt: model.Receipt{
				Retailer:     "Store5",
				Total:        15.00,
				Items:        []model.Item{{ShortDesc: "Description1234", Price: 3.0}, {ShortDesc: "Item1", Price: 2.0}},
				PurchaseDate: "2024-12-06",
				PurchaseTime: "13:00",
			},
			want:    87,
			wantErr: false,
		},
		{
			name: "Valid purchase time between 2 PM and 3 PM",
			receipt: model.Receipt{
				Retailer:     "Store6",
				Total:        50.00,
				Items:        []model.Item{{ShortDesc: "Item1", Price: 25.0}, {ShortDesc: "Item2", Price: 25.0}},
				PurchaseDate: "2024-12-06",
				PurchaseTime: "14:30",
			},
			want:    96,
			wantErr: false,
		},
		{
			name: "Valid purchase time not between 2 PM and 3 PM",
			receipt: model.Receipt{
				Retailer:     "Store7",
				Total:        30.00,
				Items:        []model.Item{{ShortDesc: "Item1", Price: 15.0}, {ShortDesc: "Item2", Price: 15.0}},
				PurchaseDate: "2024-12-06",
				PurchaseTime: "16:00",
			},
			want:    86,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CalculatePoints(tt.receipt)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CalculatePoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.CalculatePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
