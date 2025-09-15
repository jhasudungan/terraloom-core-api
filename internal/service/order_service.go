package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/constant"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepository     *repository.OrderRepository
	productRepository   *repository.ProductRepository
	orderItemRepository *repository.OrderItemRepository
	paymentRepository   *repository.PaymentRepository
	accountRepository   *repository.AccountRepository
	idGenerator         *common.IdGenerator
}

func NewOrderService(orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository, orderItemRepository *repository.OrderItemRepository, paymentRepository *repository.PaymentRepository, accountRepository *repository.AccountRepository, idGenerator *common.IdGenerator) *OrderService {
	return &OrderService{
		productRepository:   productRepository,
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
		paymentRepository:   paymentRepository,
		accountRepository:   accountRepository,
		idGenerator:         idGenerator,
	}
}

func (os *OrderService) SubmitOrder(ctx context.Context, submitOrderRequest model.SubmitOrderRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	// Input validation
	err := os.validateOrderRequest(submitOrderRequest)

	if err != nil {
		return response, err
	}

	// Use GORM transaction with callback for automatic rollback
	err = os.orderRepository.GetDB().Transaction(func(tx *gorm.DB) error {

		// Create transaction-scoped repositories
		orderRepo := repository.NewOrderRepository(tx)
		productRepo := repository.NewProductRepository(tx)
		orderItemRepo := repository.NewOrderItemRepository(tx)
		paymentRepo := repository.NewPaymentRepository(tx)
		accountRepo := repository.NewAccountRepository(tx)

		// Generate order reference
		newOrderReference, err := os.idGenerator.GenerateCommonID("ORDER")
		if err != nil {
			return err
		}

		// check for account
		account, err := accountRepo.FindByUsername(ctx, submitOrderRequest.AccountUsername)

		if err != nil {
			return err
		}

		if !account.IsActive {
			err := errors.New("account inactive")
			logrus.Error(err)
			return common.NewError(err, common.ErrAccessDenied)
		}

		newOrder := entity.Order{
			OrderReference:  newOrderReference,
			Status:          constant.OrderStatusPendingPayment,
			DeliveryAddress: submitOrderRequest.DeliveryAddress,
			OrderDate:       time.Now(),
			CreatedAt:       time.Now(),
			CreatedBy:       account.Username, // TODO: Get from context/JWT
			UpdatedAt:       time.Now(),
			UpdatedBy:       account.Username,
			AccountUsername: account.Username,
		}

		// Pre-validate all products exist and are available
		err = os.validateRequestProducts(ctx, productRepo, submitOrderRequest.OrderItems)

		if err != nil {
			return err
		}

		grandTotalOrder := int64(0)
		var orderItems []entity.OrderItem
		var usedProducts []entity.Product

		// Process each order item
		for _, orderItemRequest := range submitOrderRequest.OrderItems {

			// Find and lock product to prevent race conditions
			product, err := productRepo.FindByID(ctx, orderItemRequest.ProductId)

			if err != nil {
				return err
			}

			// Validate product availability
			if !product.IsActive {
				err := fmt.Errorf("product is not active: %v", orderItemRequest.ProductId)
				logrus.Error(err)
				return common.NewError(err, common.ErrValidation)
			}

			// Check stock availability using aggregated quantities
			if product.Stock < orderItemRequest.Quantity {
				err := fmt.Errorf("insufficient stock for product: %v , requested : %v , available: %v ", product.ID, orderItemRequest.Quantity, product.Stock)
				logrus.Error(err)
				return common.NewError(err, common.ErrValidation)
			}

			// lock the stock
			product.Stock = product.Stock - orderItemRequest.Quantity
			product.UpdatedAt = time.Now()
			product.UpdatedBy = constant.SYSTEM

			usedProducts = append(usedProducts, product)

			// Create order item
			orderItem, err := os.createOrderItem(orderItemRequest, newOrder, account)

			if err != nil {
				return err
			}

			orderItems = append(orderItems, orderItem)

			// Calculate grand total
			newGrandTotal := grandTotalOrder + orderItem.Total

			if newGrandTotal < grandTotalOrder {
				err := fmt.Errorf("grand total overflow")
				logrus.Error(err)
				return common.NewError(err, common.ErrConflict)
			}

			grandTotalOrder = newGrandTotal

			// Business rule: Maximum order total
			if grandTotalOrder > 10000000000 {
				err := fmt.Errorf("order total exceeds maximum limit: %v", grandTotalOrder)
				logrus.Error(err)
				return common.NewError(err, common.ErrValidation)
			}

		}

		// Set order total and create order
		newOrder.Total = grandTotalOrder

		err = orderRepo.Create(ctx, newOrder)
		if err != nil {
			return err
		}

		for i, oi := range orderItems {
			logrus.Infof("orderItem[%d] ref=%s product=%d qty=%d", i, oi.OrderItemReference, oi.ProductID, oi.Quantity)
		}

		err = orderItemRepo.CreateBatch(ctx, orderItems, len(orderItems))

		if err != nil {
			return err
		}

		err = productRepo.BatchUpsert(ctx, usedProducts)

		if err != nil {
			return err
		}

		// Create payment with status pending
		newPayment, err := os.createPayment(newOrder, account)

		if err != nil {
			return err
		}

		err = paymentRepo.Create(ctx, newPayment)

		if err != nil {
			return err
		}

		// Prepare response data
		responseData := model.SubmitOrderResponseData{
			OrderReference: newOrderReference,
			Total:          grandTotalOrder,
			OrderDate:      newOrder.OrderDate,
			OrderStatus:    newOrder.Status,
		}

		response.ResponseCode = constant.SuccessCode
		response.ResponseMessage = constant.SuccessMessage
		response.Data = responseData

		logrus.Info("Order created successfully:", newOrderReference, "total:", grandTotalOrder)
		return nil
	})

	// Handle transaction result
	if err != nil {
		return response, err
	}

	return response, nil
}

/**
	Unexported function (internal use only)
**/

/*
*

	Bussiness Rule To Prevent Abuse Of System, :
	- Order need to have order items
	- Maximum order items per order             : 100 Item
	- Maximum purchased quantity per order item : 1000
	- Maximum grand total quantity per order    : 10000

*
*/
func (os *OrderService) validateOrderRequest(request model.SubmitOrderRequest) error {

	if len(request.OrderItems) < 1 {
		err := errors.New("empty order items")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	if len(request.OrderItems) > 100 {
		err := errors.New("too many order items")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	totalQuantity := int64(0)

	for i, item := range request.OrderItems {

		// Basic validation
		if item.ProductId == 0 {
			err := errors.New("invalid product ID at index " + string(rune(i)))
			logrus.Error(err)
			return common.NewError(err, common.ErrValidation)
		}

		if item.Quantity <= 0 {
			err := errors.New("invalid quantity at index " + string(rune(i)))
			logrus.Error(err)
			return common.NewError(err, common.ErrValidation)
		}

		if item.PriceUsed <= 0 {
			err := errors.New("invalid price at index " + string(rune(i)))
			logrus.Error(err)
			return common.NewError(err, common.ErrValidation)
		}

		if item.Quantity > 1000 { // Max quantity per line item
			err := errors.New("quantity too large at index " + string(rune(i)))
			logrus.Error(err)
			return common.NewError(err, common.ErrValidation)
		}

		totalQuantity = totalQuantity + item.Quantity

	}

	if totalQuantity > 10000 {
		err := errors.New("grand total quantity too large")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	return nil
}

/*
*

	validate all the items ready to procces

*
*/
func (os *OrderService) validateRequestProducts(ctx context.Context, productRepo *repository.ProductRepository, orderItems []model.OrderItemRequest) error {

	// Collect product IDs
	productIDs := make([]int64, 0, len(orderItems))

	for _, item := range orderItems {
		productIDs = append(productIDs, item.ProductId)
	}

	// Batch fetch all products to verify they exist
	products, err := productRepo.FindMultipleByIDs(ctx, productIDs)
	if err != nil {
		logrus.Error(err)
		return common.NewError(err, common.ErrResourceNotFound)
	}

	// Verify all products were found
	if len(products) != len(productIDs) {

		foundIDs := make(map[int64]bool)
		for _, product := range products {
			foundIDs[product.ID] = true
		}

		for _, requiredID := range productIDs {
			if !foundIDs[requiredID] {
				err := fmt.Errorf("product not found: %d", requiredID)
				logrus.Error(err)
				return common.NewError(err, common.ErrValidation)
			}
		}
	}

	return nil
}

func (os *OrderService) createOrderItem(orderItemRequest model.OrderItemRequest, order entity.Order, account entity.Account) (entity.OrderItem, error) {

	newOrderItemReference, err := os.idGenerator.GenerateCommonID("OI")

	if err != nil {
		return entity.OrderItem{}, err
	}

	total := orderItemRequest.PriceUsed * orderItemRequest.Quantity

	// Prevent integer overflow
	if total < 0 || total < orderItemRequest.PriceUsed || total < orderItemRequest.Quantity {
		err := fmt.Errorf("price calculation overflow for product: %v", orderItemRequest.ProductId)
		logrus.Error(err)
		return entity.OrderItem{}, common.NewError(err, common.ErrValidation)
	}

	return entity.OrderItem{
		OrderItemReference:      newOrderItemReference,
		OrderReference:          order.OrderReference,
		ProductID:               orderItemRequest.ProductId,
		Quantity:                orderItemRequest.Quantity,
		PriceSnapshot:           orderItemRequest.PriceUsed,
		ProductNameSnapshot:     orderItemRequest.ProductName,
		ProductImageUrlSnapshot: orderItemRequest.ProductImageUrl,
		Total:                   total,
		CreatedAt:               time.Now(),
		CreatedBy:               account.Username,
		UpdatedAt:               time.Now(),
		UpdatedBy:               account.Username,
	}, nil
}

func (os *OrderService) createPayment(order entity.Order, account entity.Account) (entity.Payment, error) {

	newPaymentReference, err := os.idGenerator.GenerateCommonID("PAY")

	if err != nil {
		return entity.Payment{}, err
	}

	return entity.Payment{
		PaymentReference: newPaymentReference,
		OrderReference:   order.OrderReference,
		Status:           constant.PaymentStatusPending,
		Total:            order.Total,
		CreatedAt:        time.Now(),
		CreatedBy:        account.Username,
		UpdatedAt:        time.Now(),
		UpdatedBy:        account.Username,
	}, nil
}

func (os *OrderService) CancelOrder(ctx context.Context, request model.CancelOrderRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	order, err := os.orderRepository.FindByIDWithItems(ctx, request.OrderReference)

	if err != nil {
		return response, err
	}

	account, err := os.accountRepository.FindByUsername(ctx, request.AccountUsername)

	if err != nil {
		return response, err
	}

	// Order Processed can't be undone
	if order.Status == constant.OrderStatusProcessed {
		err := errors.New("order status already final")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrConflict)
	}

	payment, err := os.paymentRepository.FindByOrderReference(ctx, request.OrderReference)

	if err != nil {
		return response, err
	}

	err = os.orderRepository.GetDB().Transaction(func(tx *gorm.DB) error {

		orderRepo := repository.NewOrderRepository(tx)
		paymentRepo := repository.NewPaymentRepository(tx)
		productRepo := repository.NewProductRepository(tx)

		if order.Status == constant.OrderStatusPendingPayment {
			payment.Status = constant.PaymentStatusCancelled
		}

		if order.Status == constant.OrderStatusPaymentReceived {
			payment.Status = constant.PaymentStatusRefunded
		}

		order.Status = constant.OrderStatusCancelled
		order.UpdatedBy = account.Username
		order.UpdatedAt = time.Now()
		payment.UpdatedBy = account.Username
		payment.UpdatedAt = time.Now()

		err = orderRepo.Update(ctx, order)

		if err != nil {
			return err
		}

		err = paymentRepo.Update(ctx, payment)

		if err != nil {
			return err
		}

		// Return product stock
		productIDs := make([]int64, 0, len(order.OrderItems))

		for _, item := range order.OrderItems {
			productIDs = append(productIDs, item.ProductID)
		}

		// find the returned
		usedProducts, err := productRepo.FindMultipleByIDs(ctx, productIDs)

		if err != nil {
			return err
		}

		// Create a map for O(1) lookup
		productMap := make(map[int64]entity.Product)

		for i := range usedProducts {
			productMap[usedProducts[i].ID] = usedProducts[i]
		}

		// Update stock quantities
		productsToUpdate := make([]entity.Product, 0, len(usedProducts))

		for _, orderItem := range order.OrderItems {
			if product, exists := productMap[orderItem.ProductID]; exists {
				product.Stock += orderItem.Quantity
				product.UpdatedAt = time.Now()
				product.UpdatedBy = constant.SYSTEM
				productsToUpdate = append(productsToUpdate, product)
			}
		}

		// Batch update all products

		for _, productToUpdate := range productsToUpdate {
			logrus.Info(productToUpdate)
		}

		if len(productsToUpdate) > 0 {
			err = productRepo.BatchUpsert(ctx, productsToUpdate)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	responseData := model.CancelOrderResponseData{
		OrderReference: order.OrderReference,
		OrderDate:      order.OrderDate,
		OrderStatus:    order.Status,
		PaymentStatus:  payment.Status,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil

}

func (os *OrderService) GetAccountOrders(ctx context.Context, request model.GetAccountOrdersRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	// validate pagination
	if request.IsPaginate {
		if request.Page < 1 {
			err := errors.New("page must be greater than 0")
			logrus.Error(err)
			return response, common.NewError(err, common.ErrAccessDenied)
		}
		if request.PerPage < 1 || request.PerPage > 100 {
			err := errors.New("page must be greater than 0")
			logrus.Error(err)
			return response, common.NewError(err, common.ErrAccessDenied)
		}
	}

	// Build param for repository
	filter := model.OrderFilter{
		OrderReference: request.OrderReference,
	}

	paginationParams := model.PaginationParams{
		IsPaginate: request.IsPaginate,
		Page:       request.Page,
		PerPage:    request.PerPage,
	}

	orders, totalData, err := os.orderRepository.FindWithAccountAndFilters(ctx, request.AccountUserame, filter, paginationParams)

	if err != nil {
		return response, err
	}

	ordersDTO := make([]model.OrderDTO, len(orders))

	for i, order := range orders {
		ordersDTO[i] = model.OrderDTO{
			OrderReference: order.OrderReference,
			Status:         order.Status,
			OrderDate:      order.OrderDate,
			Total:          order.Total}
	}

	metadata := model.MetadataDTO{}
	metadata.Page = request.Page
	metadata.TotalData = totalData

	totalPage := 1
	if request.IsPaginate && request.PerPage > 0 {
		totalPage = int(math.Ceil(float64(totalData) / float64(request.PerPage)))
	}

	metadata.TotalPage = totalPage
	metadata.PerPage = request.PerPage

	responseData := model.GetAccountOrdersResponseData{
		Orders:   ordersDTO,
		Metadata: metadata,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, err

}

func (os *OrderService) GetOrderDetail(ctx context.Context, request model.GetOrderDetailRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	order, err := os.orderRepository.FindByIDWithItems(ctx, request.OrderReference)

	if err != nil {
		return response, err
	}

	payment, err := os.paymentRepository.FindByOrderReference(ctx, request.OrderReference)

	if err != nil {
		return response, err
	}

	var orderItemsDTO []model.OrderItemDTO

	for _, orderItem := range order.OrderItems {

		orderItemProduct := model.OrderItemProductDTO{
			ID:       orderItem.ProductID,
			Name:     orderItem.ProductNameSnapshot,
			ImageUrl: orderItem.ProductImageUrlSnapshot,
		}

		orderItemDTO := model.OrderItemDTO{
			OrderItemReference: orderItem.OrderItemReference,
			Quantity:           orderItem.Quantity,
			Price:              orderItem.PriceSnapshot,
			Total:              orderItem.Total,
			Product:            orderItemProduct,
		}

		orderItemsDTO = append(orderItemsDTO, orderItemDTO)
	}

	paymentDTO := model.PaymentDTO{
		PaymentReference: payment.PaymentReference,
		Status:           payment.Status,
		Total:            payment.Total,
		CardHolderName:   payment.CardHolderName,
		CardNumber:       payment.CardNumber,
	}

	orderDTO := model.OrderWithPaymentAndItemDTO{
		OrderReference:  order.OrderReference,
		OrderDate:       order.OrderDate,
		DeliveryAddress: order.DeliveryAddress,
		Status:          order.Status,
		Total:           order.Total,
		Payment:         paymentDTO,
		OrderItems:      orderItemsDTO,
	}

	responseData := model.GetOrderDetailReponseData{
		Order: orderDTO,
	}

	response.ResponseMessage = constant.SuccessMessage
	response.ResponseCode = constant.SuccessCode
	response.Data = responseData

	return response, nil
}
