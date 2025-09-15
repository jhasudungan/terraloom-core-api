package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	accountService *service.AccountService
	orderService   *service.OrderService
	errorHandler   *ErrorHandler
}

func NewAccountHandler(
	accountService *service.AccountService,
	orderService *service.OrderService,
	errorHandler *ErrorHandler) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		orderService:   orderService,
		errorHandler:   errorHandler,
	}
}

func (ah *AccountHandler) Register(ctx *gin.Context) {

	request := model.RegisterRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	response, err := ah.accountService.Register(ctx, request)

	if err != nil {
		ah.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (ah *AccountHandler) Login(ctx *gin.Context) {

	request := model.LoginRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	response, err := ah.accountService.Login(ctx, request)

	if err != nil {
		ah.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (ah *AccountHandler) GetAccountDetail(ctx *gin.Context) {

	request := model.GetAccountDetailRequest{}
	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.Username = username.(string)

	response, err := ah.accountService.GetAccountDetail(ctx, request)

	if err != nil {
		ah.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (oh *AccountHandler) GetAccountOrders(ctx *gin.Context) {

	request := model.GetAccountOrdersRequest{}

	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.AccountUserame = username.(string)
	request.OrderReference = ctx.Query("orderReference")

	isPaginate, err := strconv.ParseBool(ctx.Query("isPaginate"))

	if err != nil {
		logrus.Error(err)
		oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.IsPaginate = isPaginate

	if isPaginate {

		page, err := strconv.Atoi(ctx.Query("page"))

		if err != nil {
			logrus.Error(err)
			oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
			return
		}

		request.Page = page

		perPage, err := strconv.Atoi(ctx.Query("perPage"))

		if err != nil {
			logrus.Error(err)
			oh.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
			return
		}

		request.PerPage = perPage

	} else {
		request.Page = 1
		request.PerPage = 1
	}

	response, err := oh.orderService.GetAccountOrders(ctx, request)

	if err != nil {
		oh.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)

}

func (ah *AccountHandler) UpdateAccount(ctx *gin.Context) {

	request := model.UpdateAccountRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.Username = username.(string)
	response, err := ah.accountService.UpdateAccount(ctx, request)

	if err != nil {
		ah.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (ah *AccountHandler) UpdatePassword(ctx *gin.Context) {

	request := model.UpdatePasswordRequest{}
	err := ctx.ShouldBind(&request)

	if err != nil {
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	username, exists := ctx.Get("username")

	if !exists {
		err := errors.New("missing required data")
		logrus.Error(err)
		ah.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.Username = username.(string)
	response, err := ah.accountService.UpdatePassword(ctx, request)

	if err != nil {
		ah.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}
