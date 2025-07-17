package handler

import (
	"net/http"
	"strconv"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/gin-gonic/gin"
)

// CreateRole godoc
// @Router      /role [post]
// @Summary     Create a role
// @Description Create a new role
// @Security    BearerAuth
// @Tags        role
// @Accept      json
// @Produce     json
// @Param       body body models.Role true "Role"
// @Success     201 {object} models.Role
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) CreateRole(ctx *gin.Context) {
	body := models.Role{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.Log.Error("error while decoding request body", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	body.CreatedBy = ctx.GetString("user_id")

	role, err := h.Storage.Role().Create(ctx.Request.Context(), body)
	if err != nil {
		h.Log.Error("error while creating role", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

// GetAllRoles godoc
// @Router      /role/list [get]
// @Summary     Get all roles
// @Description Get all roles
// @Security    BearerAuth
// @Tags        role
// @Accept      json
// @Produce     json
// @Param       search query string false "Search by name"
// @Param       page query int false "Page number" default(1)
// @Param       limit query int false "Number of items per page" default(10)
// @Success     200 {object} models.GetListRoleResponse
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) GetListRoles(ctx *gin.Context) {
	request := models.GetListRequest{}

	request.Search = ctx.Query("search")

	page, err := strconv.ParseUint(ctx.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "invalid page number"})
		return
	}

	limit, err := strconv.ParseUint(ctx.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "invalid limit number"})
		return
	}

	request.Page = page
	request.Limit = limit

	response, err := h.Storage.Role().GetList(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "failed to get roles list",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    response,
		"message": "Roles retrieved successfully",
	})
}

// GetSingleRole godoc
// @Router      /role/{id} [get]
// @Summary     Get a single role
// @Description Get a single role by ID or name
// @Security    BearerAuth
// @Tags        role
// @Accept      json
// @Produce     json
// @Param       id query string false "Role ID"
// @Success     200 {object} models.Role
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) GetSingleRole(ctx *gin.Context) {
	id := ctx.Query("id")

	request := models.ID{
		Id: id,
	}

	role, err := h.Storage.Role().GetSingle(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":  http.StatusInternalServerError,
			"error": "failed to get role",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    role,
		"message": "Role retrieved successfully",
	})
}

// UpdateRole godoc
// @Router      /role/{id} [put]
// @Summary     Update a role
// @Description Update a role by ID
// @Security    BearerAuth
// @Tags        role
// @Accept      json
// @Produce     json
// @Param       id path string true "Role ID"
// @Param       body body models.Role true "Role"
// @Success     200 {object} models.Role
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) UpdateRole(ctx *gin.Context) {

	id := ctx.Param("id")
	body := models.Role{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.Log.Error("error while decoding request body", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	body.Id = id

	role, err := h.Storage.Role().Update(ctx.Request.Context(), body)
	if err != nil {
		h.Log.Error("error while updating role", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}
