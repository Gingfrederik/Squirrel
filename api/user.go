package api

import (
	"fileserver/config"
	"fileserver/types"
	"fileserver/user"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func (h *Handler) login(c *gin.Context) {
	cfg := config.GetInstance()
	userM := user.GetInstance()

	userData := &types.User{}
	err := c.BindJSON(userData)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := userM.Login(userData)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := generateJWTToken(dbUser, cfg.SecretKey)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	expiresIn := time.Hour * time.Duration(1)

	token := types.Token{
		AccessToken: jwtToken,
		TokenType:   "bearer",
		ExpiresIN:   int(expiresIn.Seconds()),
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) register(c *gin.Context) {
	userM := user.GetInstance()
	userData := &types.User{}
	err := c.BindJSON(userData)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = userM.Register(userData)
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) getAllUser(c *gin.Context) {
	userM := user.GetInstance()

	users, err := userM.GetAllUser()
	if err != nil {
		abortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := types.Response{
		Status:  0,
		Message: "get all users",
		Data:    users,
	}

	c.JSON(http.StatusOK, res)
}

func generateJWTToken(user *types.User, secretKey string) (token string, err error) {
	nowTime := time.Now()
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = user.ID.String()
	claims["name"] = user.Name
	claims["exp"] = nowTime.Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = nowTime.Unix()
	jwtToken.Claims = claims

	token, err = jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return
	}

	return
}
