package models

import "time"

type Otp struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetSingleOTP struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Email  string `json:"email"`
}

type OtpRequest struct {
	Email string `json:"email"`
}

type OtpConfirmRequest struct {
	OtpID string `json:"otp_id"`
	Code  string `json:"code"`
}

type RegisterRequest struct {
	OtpConfirmationToken string `json:"otp_confirmation_token"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	Name                 string `json:"name"`
}
