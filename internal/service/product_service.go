package service

import (
	"context"
	"errors"
	"math"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/constant"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/repository"
	"github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

func (ps *ProductService) GetProducts(ctx context.Context, request model.GetProductsRequest) (model.GeneralResponse, error) {

	// response
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
	filter := model.ProductFilter{
		Name:     request.Name,
		IsActive: request.IsActive,
	}

	paginationParams := model.PaginationParams{
		IsPaginate: request.IsPaginate,
		Page:       request.Page,
		PerPage:    request.PerPage,
	}

	products, totalData, err := ps.productRepository.FindWithFilters(ctx, filter, paginationParams)

	if err != nil {
		return response, err
	}

	productsDTO := make([]model.ProductDTO, len(products))

	for i, product := range products {
		productsDTO[i] = model.ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			CategoryID:  product.CategoryID,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
			Description: product.Description,
			ImageUrl:    product.ImageUrl,
		}
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

	responseData := model.GetProductsResponseData{
		Products: productsDTO,
		Metadata: metadata,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, err
}

func (ps *ProductService) GetProductDetail(ctx context.Context, request model.GetProductDetailRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	product, err := ps.productRepository.FindByID(ctx, request.ID)

	if err != nil {
		return response, err
	}

	productsDTO := model.ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		Description: product.Description,
		ImageUrl:    product.ImageUrl,
	}

	responseData := model.GetProductDetailResponseData{
		Product: productsDTO,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil

}
