package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	errorHandler   *ErrorHandler
}

func NewPaymentHandler(
	paymentService *service.PaymentService,
	errorHandler *ErrorHandler) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		errorHandler:   errorHandler,
	}
}

func (ph *PaymentHandler) SubmitPayment(ctx *gin.Context) {

	// Parse query parameters
	request := model.SubmitPaymentRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	response, err := ph.paymentService.SubmitPayment(ctx, request)

	if err != nil {
		ph.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}
