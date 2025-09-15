package model

type GetProductsRequest struct {
	Name       string
	IsActive   bool
	Page       int
	PerPage    int
	IsPaginate bool
}

type GetProductDetailRequest struct {
	ID int64
}

type OrderItemRequest struct {
	ProductId       int64  `json:"productId"`
	PriceUsed       int64  `json:"priceUsed"`
	Quantity        int64  `json:"quantity"`
	ProductName     string `json:"productName"`
	ProductImageUrl string `json:"productImageUrl"`
}

type SubmitOrderRequest struct {
	AccountUsername string
	DeliveryAddress string             `json:"deliveryAddress"`
	OrderItems      []OrderItemRequest `json:"orderItems"`
}

type CancelOrderRequest struct {
	OrderReference  string `json:"orderReference"`
	AccountUsername string
}

type SubmitPaymentRequest struct {
	OrderReference string `json:"orderReference"`
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	Status         string `json:"status"`
}

type GetOrderDetailRequest struct {
	OrderReference string `json:"orderReference"`
}

type RegisterRequest struct {
	Username          string `json:"username"`
	DiplayName        string `json:"displayName"`
	Email             string `json:"email"`
	LoginPassword     string `json:"loginPassword"`
	RegisteredAddress string `json:"registeredAddress"`
}

type UpdateAccountRequest struct {
	DiplayName        string `json:"displayName"`
	Email             string `json:"email"`
	RegisteredAddress string `json:"registeredAddress"`
	Username          string `json:"username"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Username    string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetAccountDetailRequest struct {
	Username string
}

type GetAccountOrdersRequest struct {
	AccountUserame string
	OrderReference string
	Page           int
	PerPage        int
	IsPaginate     bool
}
