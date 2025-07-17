package handler

import (
	"database/sql"
	"net/http"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateCustomer godoc
// @Router      /sysuser [POST]
// @Summary     Create a sysuser
// @Description Create a new sysuser
// @Security    BearerAuth
// @Tags        sysuser
// @Accept      json
// @Produce 	json
// @Param 		sysuser body models.SysUser true "sysuser"
// @Success 	201  {object}  models.SuccessResponse
// @Failure		400  {object}  models.ErrorResponse
func (h *Handler) CreateSysuser(ctx *gin.Context) {
	req := models.SysUser{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Step 1: Check if sysuser already exists
	existingUser, err := h.Storage.SysUser().GetByEmailAndStatus(ctx, req.Email, []string{"active", "blocked"})
	if err != nil && err != sql.ErrNoRows {
		h.Log.Error("failed to get sysuser", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "internal error",
			"description": err.Error(),
		})
		return
	}
	if err == nil && existingUser.Id != "" {
		ctx.JSON(http.StatusConflict, gin.H{"error": "sysuser already exists"})
		return
	}

	// Step 2: Check each role exists
	for _, roleID := range req.Roles {
		exists, err := h.Storage.Role().ExistsByIDAndStatus(ctx, roleID, "active")
		if err != nil {
			h.Log.Error("failed to check role", logger.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "role not found", "role_id": roleID})
			return
		}
	}

	// Step 3: Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error("failed to hash password", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Step 4: Create sysuser
	sysuser := models.SysUser{
		Id:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   "active",
	}

	if err := h.Storage.SysUser().Create(ctx, sysuser); err != nil {
		h.Log.Error("failed to create sysuser", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sysuser"})
		return
	}

	// Step 5: Attach roles
	for _, roleID := range req.Roles {
		if err := h.Storage.SysUser().AttachRole(ctx, sysuser.Id, roleID); err != nil {
			h.Log.Error("failed to attach role", logger.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role", "role_id": roleID})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": sysuser.Id})
}
