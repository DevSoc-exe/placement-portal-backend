package models

import (
	"database/sql"
	"time"
)

// ! Phone Number to be included in the User struct
type User struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	PhoneNumber       string         `json:"phone_number"`
	Email             string         `json:"email"`
	Gender            string         `json:"gender"`
	Otp               sql.NullString `json:"otp"`
	RollNumber        string         `json:"rollnum"`
	YearOfAdmission   int            `json:"year_of_admission"`
	Branch            string         `json:"branch"`
	StudentType       string         `json:"student_type"`
	RefreshToken      string         `json:"refresh_token"`
	IsVerified        bool           `json:"is_verified"`
	VerificationToken sql.NullString `json:"verification_token"`
	Role              string         `json:"role"`
	IsOnboarded       bool           `json:"isOnboarded"`
}

type UserResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender"`
	RollNumber      string `json:"rollnum"`
	YearOfAdmission int    `json:"year_of_admission"`
	Branch          string `json:"branch"`
	StudentType     string `json:"student_type"`
	IsOnboarded     bool   `json:"isOnboarded"`
	Role			string `json:"role"`
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
	Gender          string `json:"gender" binding:"required"`
	RollNum         string `json:"rollnum" binding:"required"`
	YearOfAdmission int    `json:"year_of_admission" binding:"required"`
	Branch          string `json:"branch"`
	StudentType     string `json:"student_type"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
}

type LoginResponse struct {
	Email string `json:"email"`
}
