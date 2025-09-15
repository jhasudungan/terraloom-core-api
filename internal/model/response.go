package model

import "time"

type GeneralResponse struct {
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data"`
}

type ErrorResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Detail          string `json:"detail"`
}

type MetadataDTO struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"perPage"`
	TotalData int64 `json:"totalData"`
	TotalPage int   `json:"totalPage"`
}

type GetProductsResponseData struct {
	Products []ProductDTO `json:"products"`
	Metadata MetadataDTO  `json:"metadata"`
}

type GetProductDetailResponseData struct {
	Product ProductDTO `json:"product"`
}

type SubmitOrderResponseData struct {
	OrderReference string    `json:"orderReference"`
	OrderDate      time.Time `json:"orderDate"`
	OrderStatus    string    `json:"orderStatus"`
	Total          int64     `json:"total"`
}

type CancelOrderResponseData struct {
	OrderReference string    `json:"orderReference"`
	OrderDate      time.Time `json:"orderDate"`
	OrderStatus    string    `json:"orderStatus"`
	PaymentStatus  string    `json:"paymentStatus"`
}

type RegisterResponseData struct {
	Account AccountDTO `json:"account"`
}

type LoginRepsonseData struct {
	Token   TokenDTO   `json:"token"`
	Account AccountDTO `json:"account"`
}

type UpdateAccountResponseData struct {
	Account AccountDTO `json:"account"`
}

type GetAccountDetailResponseData struct {
	Account AccountDTO `json:"account"`
}

type GetAccountOrdersResponseData struct {
	Orders   []OrderDTO  `json:"orders"`
	Metadata MetadataDTO `json:"metadata"`
}

type GetOrderDetailReponseData struct {
	Order OrderWithPaymentAndItemDTO `json:"order"`
}
