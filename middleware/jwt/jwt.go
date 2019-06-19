package jwt

import (
	"fileserver/config"
	"fileserver/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetInstance()
		r := c.Request
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			c.Next()
			return
		}
		splitToken := strings.Split(reqToken, "Bearer ")
		tokenString := splitToken[1]
		if tokenString == "" {
			res := types.Response{
				Status:  -1,
				Message: "unauthorization",
			}
			c.JSON(http.StatusForbidden, res)
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.SecretKey), nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					res := types.Response{
						Status:  -2,
						Message: "token expired",
					}
					c.JSON(http.StatusForbidden, res)
					c.Abort()
					return
				} else {
					res := types.Response{
						Status:  -3,
						Message: "token invaild",
					}
					c.JSON(http.StatusForbidden, res)
					c.Abort()
					return
				}
			}
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", token.Claims)
		}

		c.Next()
	}
}
