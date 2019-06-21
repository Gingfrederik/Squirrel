package api

import (
	"fileserver/auth"
	"fileserver/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllPolicy(c *gin.Context) {
	authM := auth.GetInstance()
	policy := authM.GetPolicy()

	res := types.Response{
		Status:  0,
		Message: "get all policy",
		Data:    policy,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) addPolicy(c *gin.Context) {
	authM := auth.GetInstance()

	policy := &types.Policy{}
	err := c.BindJSON(&policy)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
	}

	result := authM.AddPermissionForUser(policy.Role, policy.Path, policy.Method)
	if result {
		authM.SaveAndReloadPolicy()
	}

	res := types.Response{
		Status:  0,
		Message: "Success add policy",
		Data:    result,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) delPolicy(c *gin.Context) {
	authM := auth.GetInstance()

	policy := &types.Policy{}
	err := c.BindJSON(&policy)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
	}

	result := authM.DeletePermissionForUser(policy.Role, policy.Path, policy.Method)
	if result {
		authM.SaveAndReloadPolicy()
	}

	res := types.Response{
		Status:  0,
		Message: "Success delete policy",
		Data:    result,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllRole(c *gin.Context) {
	authM := auth.GetInstance()
	roles := authM.GetAllRoles()

	res := types.Response{
		Status:  0,
		Message: "get all roles",
		Data:    roles,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllUserRole(c *gin.Context) {
	authM := auth.GetInstance()
	roles := authM.GetAllRoles()

	data := []*types.RoleUser{}
	for _, role := range roles {
		users := authM.GetUsersForRole(role)
		roleUser := &types.RoleUser{
			Role:  role,
			Users: users,
		}

		data = append(data, roleUser)
	}

	res := types.Response{
		Status:  0,
		Message: "get all the users has role",
		Data:    data,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) addUserRole(c *gin.Context) {
	authM := auth.GetInstance()

	roleUser := &types.RoleUser{}
	err := c.BindJSON(&roleUser)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
	}

	result := authM.AddRoleForUser(roleUser.Users[0], roleUser.Role)
	if result {
		authM.SaveAndReloadPolicy()
	}

	res := types.Response{
		Status:  0,
		Message: "Success add role for user",
		Data:    result,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) delUserRole(c *gin.Context) {
	authM := auth.GetInstance()

	roleUser := &types.RoleUser{}
	err := c.BindJSON(&roleUser)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
	}

	result := authM.DeleteRoleForUser(roleUser.Users[0], roleUser.Role)
	if result {
		authM.SaveAndReloadPolicy()
	}

	res := types.Response{
		Status:  0,
		Message: "Success delete role for user",
		Data:    result,
	}

	c.JSON(http.StatusOK, res)
}
