package handler

import (
	"net/http"
	"time"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/config"
	"github.com/Akrom0181/Auth/pkg/etc"
	"github.com/Akrom0181/Auth/pkg/hash"
	"github.com/Akrom0181/Auth/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login godoc
// @Router      /auth/login [post]
// @Summary     Login
// @Description Login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body models.LoginRequest true "User"
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) Login(ctx *gin.Context) {

	body := models.LoginRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var userId string

	switch body.UserType {
	case "sysuser":
		sysuser, err := h.Storage.SysUser().GetSingle(ctx, models.GetSingleSysUser{
			Email: body.Email,
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials 1"})
			return
		}

		if !hash.CheckPasswordHash(body.Password, sysuser.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials 2"})
			return
		}
		userId = sysuser.Id

	case "user":
		user, err := h.Storage.User().GetSingle(ctx, models.UserSingleRequest{
			Email: body.Email,
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials 1"})
			return
		}
		if !hash.CheckPasswordHash(body.Password, user.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials 2"})
			return
		}

		userId = user.Id

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_type"})
		return
	}

	jwtFields := map[string]interface{}{
		"user_id":   userId,
		"user_type": body.UserType,
		"exp":       time.Now().Add(time.Hour * 7 * 24).Unix(),
	}

	token, err := jwt.GenerateJWT(jwtFields, h.Config.JWTSecret)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

// SendOTP godoc
// @Router      /auth/send-otp [post]
// @Summary     Send OTP
// @Description Send OTP
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body models.OtpRequest true "User"
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) SendOTP(ctx *gin.Context) {
	body := models.OtpRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	otpCode := etc.GenerateOTP(6)
	otpID := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Minute)

	otp := models.Otp{
		Id:        otpID,
		Email:     body.Email,
		Code:      otpCode,
		Status:    "unconfirmed",
		ExpiresAt: expiresAt,
	}

	if _, err := h.Storage.Otp().Create(ctx, otp); err != nil {
		handleResponseLog(ctx, h.Log, "error while creating OTP", http.StatusInternalServerError, err.Error())
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create OTP"})
		return
	}

	emailBody, err := etc.GenerateOtpEmailBody(otpCode)
	if err != nil {
		handleResponseLog(ctx, h.Log, "error while generating email body", http.StatusInternalServerError, err.Error())
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate email body"})
		return
	}

	err = etc.SendEmail(config.SmtpHost, config.SmtpPort, config.SmtpEmail, "rdstahkoquqrnuov", body.Email, emailBody)
	if err != nil {
		handleResponseLog(ctx, h.Log, "error while sending OTP email", http.StatusInternalServerError, err.Error())
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"otp_id": otpID})
}

// ConfirmOTP godoc
// @Router      /auth/confirm-otp [post]
// @Summary     Confirm OTP
// @Description Confirm OTP
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body models.OtpConfirmRequest true "OTP Confirmation"
// @Success     200 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) ConfirmOTP(ctx *gin.Context) {
	body := models.OtpConfirmRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	otp, err := h.Storage.Otp().GetSingle(ctx, models.GetSingleOTP{
		Id:     body.OtpID,
		Status: "unconfirmed",
	})

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "OTP not found or already confirmed"})
		return
	}

	if otp.Code != body.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect OTP"})
		return
	}

	err = h.Storage.Otp().Update(ctx, models.Otp{
		Id:     otp.Id,
		Status: "confirmed",
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update OTP status"})
		return
	}

	jwtFields := map[string]interface{}{
		"otp_id": otp.Id,
		"exp":    time.Now().Add(time.Hour * 7 * 24).Unix(),
	}

	token, err := jwt.GenerateJWT(jwtFields, h.Config.JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate confirmation token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"otp_confirmation_token": token,
	})
}

// Signup godoc
// @Router      /auth/signup [post]
// @Summary     Signup
// @Description Signup
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body models.RegisterRequest true "User Registration"
// @Success     201 {object} models.SuccessResponse
// @Failure     400 {object} models.ErrorResponse
func (h *Handler) Signup(ctx *gin.Context) {
	body := models.RegisterRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	payload, err := jwt.ParseJWT(body.OtpConfirmationToken, h.Config.JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired confirmation token"})
		return
	}

	otpID, _ := payload["otp_id"].(string)
	exp, _ := payload["exp"].(float64)
	if time.Now().Unix() > int64(exp) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "otp confirmation token expired"})
		return
	}

	_, err = h.Storage.Otp().GetSingle(ctx, models.GetSingleOTP{
		Id:    otpID,
		Email: body.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "otp not found or does not match email"})
		return
	}

	user, _ := h.Storage.User().GetSingle(ctx, models.UserSingleRequest{
		Email:  body.Email,
		Status: "active",
	})
	if user.Id != "" {
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := hash.HashPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}

	newUser := models.User{
		Status:   "active",
		Name:     body.Name,
		Email:    body.Email,
		Password: hashedPassword,
	}

	if _, err := h.Storage.User().Create(ctx, newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
