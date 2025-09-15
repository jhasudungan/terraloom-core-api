package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jhasudungan/terraloom-core-api/internal/common"
	"github.com/jhasudungan/terraloom-core-api/internal/handler"
	"github.com/jhasudungan/terraloom-core-api/internal/service"
	"github.com/sirupsen/logrus"
)

func NewAuthMiddleware(jwtService *service.JwtService, errorHandler *handler.ErrorHandler) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			err := errors.New("authorization header required")
			logrus.Error(err)
			errorHandler.Handle(c, common.NewError(err, common.ErrAccessDenied))
			c.Abort()
			return
		}

		// Must be "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			err := errors.New("invalid Authorization header")
			logrus.Error(err)
			errorHandler.Handle(c, common.NewError(err, common.ErrAccessDenied))
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Parse and validate token
		token, err := jwtService.ParseJWT(tokenStr)
		if err != nil || !token.Valid {
			err := errors.New("invalid token")
			logrus.Error(err)
			errorHandler.Handle(c, common.NewError(err, common.ErrAccessDenied))
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			err := errors.New("invalid token claims")
			logrus.Error(err)
			errorHandler.Handle(c, common.NewError(err, common.ErrAccessDenied))
			c.Abort()
			return
		}

		// Put `sub` into Gin context
		if sub, ok := claims["sub"].(string); ok {
			c.Set("username", sub)
		} else {
			err := errors.New("invalid token claims")
			logrus.Error(err)
			errorHandler.Handle(c, common.NewError(err, common.ErrAccessDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}
