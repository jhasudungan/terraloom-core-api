package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/sirupsen/logrus"
)

type JwtService struct {
	Secret string
}

func NewJwtService(secret string) *JwtService {
	return &JwtService{
		Secret: secret,
	}
}

func (j *JwtService) GenerateJWT(username string, expiry int64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": expiry, // 24 hours
	})

	tokenString, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		logrus.Error(err)
		return "", common.NewError(err, common.ErrAuthFailed)
	}

	return tokenString, nil
}

func (j *JwtService) ParseJWT(tokenStr string) (*jwt.Token, error) {

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Error(jwt.ErrSignatureInvalid)
			return nil, common.NewError(jwt.ErrSignatureInvalid, common.ErrAuthFailed)
		}
		return []byte(j.Secret), nil
	})
}
