package model

import "time"

type ProductDTO struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"categoryId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int64  `json:"stock"`
	Price       int64  `json:"price"`
	ImageUrl    string `json:"imageUrl"`
	IsActive    bool   `json:"isActive"`
}

type OrderDTO struct {
	OrderReference string    `json:"orderReference"`
	OrderDate      time.Time `json:"orderDate"`
	Status         string    `json:"status"`
	Total          int64     `json:"total"`
}

type PaymentDTO struct {
	PaymentReference string `json:"paymentReference"`
	Status           string `json:"status"`
	Total            int64  `json:"total"`
	CardHolderName   string `json:"cardHolderName"`
	CardNumber       string `json:"cardNumber"`
}

type OrderItemProductDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type OrderItemDTO struct {
	OrderItemReference string              `json:"orderItemReference"`
	Quantity           int64               `json:"quantity"`
	Price              int64               `json:"price"`
	Total              int64               `json:"total"`
	Product            OrderItemProductDTO `json:"product"`
}

type AccountDTO struct {
	ID                int64  `json:"id"`
	Username          string `json:"username"`
	DisplayName       string `json:"displayName"`
	Email             string `json:"email"`
	RegisteredAddress string `json:"registeredAddress"`
	IsActive          bool   `json:"isActive"`
}

type TokenDTO struct {
	Token     string `json:"token"`
	Expiry    int64  `json:"expiry"`
	ExpiredAt string `json:"expiredAt"`
}
