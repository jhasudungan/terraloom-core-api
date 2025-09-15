package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderService *service.OrderService
	errorHandler *ErrorHandler
}

func NewOrderHandler(
	orderService *service.OrderService,
	errorHandler *ErrorHandler) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		errorHandler: errorHandler,
	}
}

func (oh *OrderHandler) SubmitOrder(ctx *gin.Context) {

	request := model.SubmitOrderRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.AccountUsername = username.(string)

	response, err := oh.orderService.SubmitOrder(ctx, request)

	if err != nil {
		oh.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (oh *OrderHandler) CancelOrder(ctx *gin.Context) {

	request := model.CancelOrderRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.AccountUsername = username.(string)

	response, err := oh.orderService.CancelOrder(ctx, request)

	if err != nil {
		oh.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (oh *OrderHandler) GetOrderDetail(ctx *gin.Context) {

	request := model.GetOrderDetailRequest{}

	orderReference := ctx.Param("orderReference")
	request.OrderReference = orderReference

	response, err := oh.orderService.GetOrderDetail(ctx, request)

	if err != nil {
		oh.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}
