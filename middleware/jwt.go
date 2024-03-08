package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"eduData/setting"
)

type JWT struct {
	JwtSecretKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		JwtSecretKey: []byte(setting.JwtKey),
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
		//找到Authorization
		tokenString, err := c.Cookie("authentication")
		if err != nil {
			err = c.Error(errors.New("jwt: missing Authorization cookie" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		}

		//验证Authorization token
		j := NewJWT()
		username, ok := c.GetPostForm("username")
		if !ok {
			err = c.Error(errors.New("jwt: missing username in form" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		}

		if token, err := j.ParseToken(tokenString); err == nil && token.Valid && token.Claims.(jwt.MapClaims)["username"].(string) == username {
			c.Next()
		} else {
			err = c.Error(errors.New("jwt: invalid Authorization cookie")).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		}
	}
}
