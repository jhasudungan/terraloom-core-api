package service

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"time"

	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/constant"
	"github.com/jhasudungan/terraloom-core-api/internal/entity"
	"github.com/jhasudungan/terraloom-core-api/internal/model"
	"github.com/jhasudungan/terraloom-core-api/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	jwtService        *JwtService
	accountRepository *repository.AccountRepository
}

func NewAccountService(jwtService *JwtService, accountRepository *repository.AccountRepository) *AccountService {
	return &AccountService{
		jwtService:        jwtService,
		accountRepository: accountRepository,
	}
}

func (a *AccountService) Register(ctx context.Context, request model.RegisterRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	err := a.validateRegisterRequest(ctx, request)

	if err != nil {
		return response, err
	}

	// encrypt password
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.LoginPassword), 12)

	if err != nil {
		logrus.Error(err)
		return response, common.NewError(err, common.ErrValidation)
	}

	newAccount := entity.Account{
		Username:          request.Username,
		DisplayName:       request.DiplayName,
		Email:             request.Email,
		LoginPassword:     string(hashed),
		CreatedBy:         request.Username,
		UpdatedBy:         request.Username,
		UpdatedAt:         time.Now(),
		CreatedAt:         time.Now(),
		IsActive:          true,
		RegisteredAddress: request.RegisteredAddress}

	err = a.accountRepository.Create(ctx, newAccount)

	if err != nil {
		return response, err
	}

	newAccount, err = a.accountRepository.FindByUsername(ctx, request.Username)

	if err != nil {
		return response, err
	}

	accountDTO := model.AccountDTO{
		ID:                newAccount.ID,
		Username:          newAccount.Username,
		DisplayName:       newAccount.DisplayName,
		Email:             newAccount.Email,
		RegisteredAddress: newAccount.RegisteredAddress,
		IsActive:          newAccount.IsActive}

	responseData := model.RegisterResponseData{
		Account: accountDTO,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil
}

func (a *AccountService) Login(ctx context.Context, request model.LoginRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	// validate password
	err := a.validatePassword(request.Password)

	if err != nil {
		return response, err
	}

	account, err := a.accountRepository.FindByUsername(ctx, request.Username)

	if err != nil {
		return response, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.LoginPassword), []byte(request.Password))

	if err != nil {
		err = errors.New("invalid password")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrAuthFailed)
	}

	if !account.IsActive {
		err = errors.New("account inactive")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrAuthFailed)
	}

	// Format as ISO 8601 string
	expiry := time.Now().Add(24 * time.Hour).Unix()
	expiredAt := time.Now().Add(24 * time.Hour)
	expiredAtString := expiredAt.Format(time.RFC3339Nano)

	token, err := a.jwtService.GenerateJWT(account.Username, expiry)

	if err != nil {
		return response, err
	}

	tokenDTO := model.TokenDTO{
		Token:     token,
		Expiry:    expiry,
		ExpiredAt: expiredAtString,
	}

	accountDTO := model.AccountDTO{
		ID:                account.ID,
		Username:          account.Username,
		DisplayName:       account.DisplayName,
		Email:             account.Email,
		RegisteredAddress: account.RegisteredAddress,
		IsActive:          account.IsActive}

	loginResponseData := model.LoginRepsonseData{
		Token:   tokenDTO,
		Account: accountDTO}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = loginResponseData

	return response, nil

}

func (a *AccountService) GetAccountDetail(ctx context.Context, request model.GetAccountDetailRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	account, err := a.accountRepository.FindByUsername(ctx, request.Username)

	if err != nil {
		return response, err
	}

	accountDTO := model.AccountDTO{
		ID:                account.ID,
		Username:          account.Username,
		DisplayName:       account.DisplayName,
		Email:             account.Email,
		RegisteredAddress: account.RegisteredAddress,
		IsActive:          account.IsActive}

	responseData := model.GetAccountDetailResponseData{
		Account: accountDTO,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil
}

func (a *AccountService) UpdateAccount(ctx context.Context, request model.UpdateAccountRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	account, err := a.accountRepository.FindByUsername(ctx, request.Username)

	if err != nil {
		return response, err
	}

	if !a.isEmailValid(request.Email) {
		err := errors.New("email is not valid")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrValidation)
	}

	account.Email = request.Email
	account.DisplayName = request.DiplayName
	account.RegisteredAddress = request.RegisteredAddress
	account.UpdatedAt = time.Now()
	account.UpdatedBy = account.Username

	err = a.accountRepository.Update(ctx, account)

	if err != nil {
		return response, err
	}

	accountDTO := model.AccountDTO{
		ID:                account.ID,
		Username:          account.Username,
		DisplayName:       account.DisplayName,
		Email:             account.Email,
		RegisteredAddress: account.RegisteredAddress,
		IsActive:          account.IsActive}

	responseData := model.UpdateAccountResponseData{
		Account: accountDTO,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil
}

func (a *AccountService) UpdatePassword(ctx context.Context, request model.UpdatePasswordRequest) (model.GeneralResponse, error) {

	response := model.GeneralResponse{}

	account, err := a.accountRepository.FindByUsername(ctx, request.Username)

	if err != nil {
		return response, err
	}

	// match old password
	err = bcrypt.CompareHashAndPassword([]byte(account.LoginPassword), []byte(request.OldPassword))

	if err != nil {
		err = errors.New("invalid password")
		logrus.Error(err)
		return response, common.NewError(err, common.ErrAuthFailed)
	}

	// Check new password validation
	err = a.validatePassword(request.NewPassword)

	if err != nil {
		return response, err
	}

	// encrypt new password
	newHashed, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 12)

	if err != nil {
		logrus.Error(err)
		return response, common.NewError(err, common.ErrValidation)
	}

	account.LoginPassword = string(newHashed)
	account.UpdatedAt = time.Now()
	account.UpdatedBy = account.Username

	err = a.accountRepository.Update(ctx, account)

	if err != nil {
		return response, err
	}

	accountDTO := model.AccountDTO{
		ID:                account.ID,
		Username:          account.Username,
		DisplayName:       account.DisplayName,
		Email:             account.Email,
		RegisteredAddress: account.RegisteredAddress,
		IsActive:          account.IsActive}

	responseData := model.UpdateAccountResponseData{
		Account: accountDTO,
	}

	response.ResponseCode = constant.SuccessCode
	response.ResponseMessage = constant.SuccessMessage
	response.Data = responseData

	return response, nil
}

func (a *AccountService) validateRegisterRequest(ctx context.Context, request model.RegisterRequest) error {

	if request.Username == "" || request.DiplayName == "" || request.Email == "" || request.LoginPassword == "" {
		err := errors.New("one/several required data is missing")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	if len(request.Username) > 100 || len(request.DiplayName) > 100 || len(request.Email) > 200 {
		err := errors.New("one/several required data is too long")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	if !a.isEmailValid(request.Email) {
		err := errors.New("email is not valid")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	err := a.validatePassword(request.LoginPassword)
	if err != nil {
		return err
	}

	// check available username and email
	isUsernameUsed, err := a.accountRepository.CheckByUsername(ctx, request.Username)

	if err != nil {
		return err
	}

	if isUsernameUsed {
		err = errors.New("username already taken")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	isEmailUsed, err := a.accountRepository.CheckByEmail(ctx, request.Email)

	if err != nil {
		return err
	}

	if isEmailUsed {
		err = errors.New("email already taken")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	return nil
}

/*
*
- Minimum length (e.g., 8 characters)
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character
*
*/
func (a *AccountService) validatePassword(password string) error {

	var err error
	if len(password) < 8 {
		err = errors.New("password must be at least 8 characters long")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#~$%^&*()+|_.,<>?/{}\-]`)

	if !upper.MatchString(password) {
		err = errors.New("password must contain at least one uppercase letter")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}
	if !lower.MatchString(password) {
		err = errors.New("password must contain at least one lowercase letter")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}
	if !number.MatchString(password) {
		err = errors.New("password must contain at least one number")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}
	if !special.MatchString(password) {
		err = errors.New("password must contain at least one special character")
		logrus.Error(err)
		return common.NewError(err, common.ErrValidation)
	}

	return nil
}

func (a *AccountService) isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
