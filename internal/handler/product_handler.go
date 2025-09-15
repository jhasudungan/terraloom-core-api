package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	productService *service.ProductService
	errorHandler   *ErrorHandler
}

func NewProductHandler(
	productService *service.ProductService,
	errorHandler *ErrorHandler) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		errorHandler:   errorHandler,
	}
}

func (ph *ProductHandler) GetProducts(ctx *gin.Context) {

	// Parse query parameters
	request := model.GetProductsRequest{}

	request.Name = ctx.Query("name")

	isActive, err := strconv.ParseBool(ctx.Query("isActive"))

	if err != nil {
		logrus.Error(err)
		ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.IsActive = isActive
	isPaginate, err := strconv.ParseBool(ctx.Query("isPaginate"))

	if err != nil {
		logrus.Error(err)
		ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.IsPaginate = isPaginate

	if isPaginate {

		page, err := strconv.Atoi(ctx.Query("page"))

		if err != nil {
			logrus.Error(err)
			ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
			return
		}

		request.Page = page

		perPage, err := strconv.Atoi(ctx.Query("perPage"))

		if err != nil {
			logrus.Error(err)
			ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
			return
		}

		request.PerPage = perPage

	} else {
		request.Page = 1
		request.PerPage = 1
	}

	response, err := ph.productService.GetProducts(ctx, request)

	if err != nil {
		ph.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}

func (ph *ProductHandler) GetProductDetail(ctx *gin.Context) {

	// Parse query parameters
	request := model.GetProductDetailRequest{}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		logrus.Error(err)
		ph.errorHandler.Handle(ctx, common.NewError(err, common.ErrValidation))
		return
	}

	request.ID = id
	response, err := ph.productService.GetProductDetail(ctx, request)

	if err != nil {
		ph.errorHandler.Handle(ctx, err)
		return
	}

	ctx.JSON(200, response)
}
