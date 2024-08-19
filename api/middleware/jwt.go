package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"eduData/api/pub"
	"eduData/bootstrap"
)

type JWT struct {
	JwtSecretKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		JwtSecretKey: []byte(bootstrap.C.JwtKey),
	}
}

func (j *JWT) CreateToken(MyClaims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims)
	return token.SignedString(j.JwtSecretKey)
}

func (j *JWT) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.JwtSecretKey, nil
	})
}

// RequireAuthJwt jwt中间件
func RequireAuthJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(err).Error(),
			})
			return
		}

		//找到Authorization
		tokenString, err := c.Cookie("authentication")
		if err != nil {
			le.WithError(err).Error("jwt: missing Authorization cookie")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(fmt.Errorf("缺少参数 Authorization cookie %w", err)).Error(),
			})
			return
		}

		//验证Authorization token
		j := NewJWT()
		if token, err := j.ParseToken(tokenString); err == nil && token.Valid && token.Claims.(jwt.MapClaims)["username"].(string) == loginForm.Username {
			c.Next()
		} else {
			le.Error("jwt: invalid Authorization cookie")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": c.Error(fmt.Errorf("无效的Authorization cookie %w", err)).Error(),
			})
			return
		}
	}
}
