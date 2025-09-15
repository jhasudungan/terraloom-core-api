package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/constant"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentService struct {
	orderRepository   *repository.OrderRepository
	paymentRepository *repository.PaymentRepository
}

func NewPaymentService(orderRepository *repository.OrderRepository, paymentRepository *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		orderRepository:   orderRepository,
		paymentRepository: paymentRepository}
}

func (ps *PaymentService) SubmitPayment(ctx context.Context, request model.SubmitPaymentRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	// find data order
	order, err := ps.orderRepository.FindByID(ctx, request.OrderReference)

	if err != nil {
		logrus.Error(err)
		return response, common.NewError(err, common.ErrResourceNotFound)
	}

	if request.Status != constant.PaymentStatusReceived && request.Status != constant.PaymentStatusCancelled {
		err := errors.New("payment status not recognized")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrValidation)
	}

	err = ps.orderRepository.GetDB().Transaction(func(tx *gorm.DB) error {

		// Create transaction-scoped repositories
		orderRepo := repository.NewOrderRepository(tx)
		paymentRepo := repository.NewPaymentRepository(tx)

		payment, err := paymentRepo.FindByOrderReference(ctx, order.OrderReference)

		if err != nil {
			return err
		}

		// For now, only payment susccess
		if request.Status == constant.PaymentStatusReceived {
			order.Status = constant.OrderStatusPaymentReceived
			payment.Status = constant.PaymentStatusReceived
			payment.CardNumber = ps.maskCard(request.CardNumber)
			payment.CardHolderName = ps.maskName(request.CardHolderName)

			err = orderRepo.Update(ctx, order)

			if err != nil {
				return err
			}

			err = paymentRepo.Update(ctx, payment)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage

	return response, nil

}

func (ps *PaymentService) maskCard(card string) string {
	// remove spaces or dashes if needed
	card = strings.ReplaceAll(card, " ", "")
	card = strings.ReplaceAll(card, "-", "")

	if len(card) <= 4 {
		return card
	}

	masked := strings.Repeat("*", len(card)-4) + card[len(card)-4:]
	return masked
}

func (ps *PaymentService) maskName(name string) string {
	parts := strings.Fields(name)
	for i, part := range parts {
		if len(part) > 1 {
			parts[i] = string(part[0]) + strings.Repeat("*", len(part)-1)
		}
	}
	return strings.Join(parts, " ")
}
