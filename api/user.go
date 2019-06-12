package api

import (
	"fileserver/config"
	"fileserver/db"
	"fileserver/types"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) login(c *gin.Context) {
	cfg := config.GetInstance()
	user := &types.User{}
	err := c.BindJSON(user)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	dbAgent := db.GetInstance()
	id, err := dbAgent.GetIDByName(user.Name)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := dbAgent.GetUserByID(id)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		abortWithError(c, http.StatusForbidden, err.Error())
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
	dbAgent := db.GetInstance()

	user := &types.User{}
	err := c.BindJSON(user)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = dbAgent.GetIDByName(user.Name)
	if err != db.ErrKeyNotExists {
		abortWithError(c, http.StatusBadRequest, "name exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user.ID = h.genID.Generate()
	user.Password = string(hash)
	user.Locked = false
	user.CreatedAt = time.Now().UTC()
	user.UpdateAt = time.Now().UTC()

	err = dbAgent.AddOrUpdateUser(user)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
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
