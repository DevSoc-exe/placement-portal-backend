package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Otp               sql.NullString `json:"otp"`
	RollNumber        string         `json:"roll_number"`
	YearOfAdmission   int            `json:"year_of_admission"`
	Branch            string         `json:"branch"`
	StudentType       string         `json:"student_type"`
	RefreshToken      string         `json:"refresh_token"`
	IsVerified        bool           `json:"is_verified"`
	VerificationToken sql.NullString `json:"verification_token"`
	Role              string         `json:"role"`
}

type OTP struct {
	Date time.Time
	Otp  int
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

type OTPRequest struct {
	Email string `json:"email" binding:"required"`
	OTP   int    `json:"otp" binding:"required"`
}

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	RollNum         string `json:"rollnum" binding:"required"`
	YearOfAdmission int    `json:"year_of_admission" binding:"required"`
	Branch          string `json:"branch"`
	StudentType     string `json:"student_type"`
}

type LoginResponse struct {
	Email string `json:"email"`
}
