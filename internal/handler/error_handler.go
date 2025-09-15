package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/constant"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Handle(context *gin.Context, err error) {

	switch {
	case errors.Is(err, common.ErrResourceNotFound):
		response := model.ErrorResponse{
			ResponseCode:    constant.ResourceNotFoundCode,
			ResponseMessage: constant.ResourceNotFoundMessage,
			Detail:          err.Error(),
		}
		context.JSON(404, response)
		return
	case errors.Is(err, common.ErrAuthFailed):
		response := model.ErrorResponse{
			ResponseCode:    constant.AuthFailedCode,
			ResponseMessage: constant.AuthFailedMessage,
			Detail:          err.Error(),
		}
		context.JSON(401, response)
		return
	case errors.Is(err, common.ErrAccessDenied):
		response := model.ErrorResponse{
			ResponseCode:    constant.AccessDeniedCode,
			ResponseMessage: constant.AccessDeniedMessage,
			Detail:          err.Error()}
		context.JSON(403, response)
		return
	case errors.Is(err, common.ErrValidation):
		response := model.ErrorResponse{
			ResponseCode:    constant.ValidationFailedCode,
			ResponseMessage: constant.ValidationFailedMessage,
			Detail:          err.Error()}
		context.JSON(400, response)
		return
	case errors.Is(err, common.ErrConflict):
		response := model.ErrorResponse{
			ResponseCode:    constant.ConflictResourceCode,
			ResponseMessage: constant.ConflictResourceMessage,
			Detail:          err.Error()}
		context.JSON(409, response)
		return
	case errors.Is(err, common.ErrDBOperation):
		response := model.ErrorResponse{
			ResponseCode:    constant.ConflictResourceCode,
			ResponseMessage: constant.ConflictResourceMessage,
			Detail:          err.Error()}
		context.JSON(409, response)
		return
	default:
		response := model.ErrorResponse{
			ResponseCode:    constant.UnexpectedErrorCode,
			ResponseMessage: constant.UnexpectedErrorMessage,
			Detail:          err.Error(),
		}
		context.JSON(500, response)
		return
	}
}
