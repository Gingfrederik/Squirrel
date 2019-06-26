package auth

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &Authorizer{enforcer: e}

	return func(c *gin.Context) {
		if !a.CheckPermission(c) {
			a.RequirePermission(c)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type Authorizer struct {
	enforcer *casbin.Enforcer
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *Authorizer) CheckPermission(c *gin.Context) bool {
	var user interface{}
	user = "0"

	claims, ok := c.Get("claims")
	if ok {
		user = claims.(jwt.MapClaims)["id"]
	}
	method := c.Request.Method
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/v1/fs") {
		path = filepath.Join("/v1/fs", c.Param("path"))
	}
	return a.enforcer.Enforce(user, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *Authorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
