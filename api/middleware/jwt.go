package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"

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
		//找到Authorization
		tokenString, err := c.Cookie("authentication")
		if err != nil {
			err = c.Error(errors.New("jwt: missing Authorization cookie" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		}

		var loginForm LoginForm
		if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
			_ = c.Error(errors.New("middleware.RequireAuthJwt()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "表单格式错误,重新登陆后重新提交",
			})
		}

		//验证Authorization token
		j := NewJWT()
		if token, err := j.ParseToken(tokenString); err == nil && token.Valid && token.Claims.(jwt.MapClaims)["username"].(string) == loginForm.Username {
			c.Next()
		} else {
			_ = c.Error(errors.New("jwt: invalid Authorization cookie" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "身份认证失败",
			})

		}
	}
}
