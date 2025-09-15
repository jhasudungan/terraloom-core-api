package model

import "time"

type OrderWithPaymentAndItemDTO struct {
	OrderReference  string         `json:"orderReference"`
	OrderDate       time.Time      `json:"orderDate"`
	DeliveryAddress string         `json:"deliveryAddress"`
	Status          string         `json:"status"`
	Total           int64          `json:"total"`
	Payment         PaymentDTO     `json:"payment"`
	OrderItems      []OrderItemDTO `json:"orderItems"`
}
