package api

import (
	"github.com/gin-gonic/gin"

	"github.com/supermicah/dionysus-admin/internal/mods/rbac/biz"
	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

// Role management for RBAC
type Role struct {
	RoleBIZ *biz.Role
}

// Query
// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Query role list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param name query string false "Display name of role"
// @Param status query string false "Status of role (disabled, enabled)"
// @Success 200 {object} util.ResponseResult{data=[]schema.Role}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}
	result, err := a.RoleBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// Get
// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Get role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult{data=schema.Role}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleBIZ.Get(ctx, util.GetInt64Param(c, "id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

// Create
// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Create role record
// @Param body body schema.RoleForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=schema.Role}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.RoleBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

// Update
// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Update role record by ID
// @Param id path string true "unique id"
// @Param body body schema.RoleForm true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.RoleBIZ.Update(ctx, util.GetInt64Param(c, "id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// Delete
// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Delete role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBIZ.Delete(ctx, util.GetInt64Param(c, "id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
