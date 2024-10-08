package models

type User struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	Branch            string  `json:"branch"`
	RollNumber        string  `json:"rollnum"`
	YearOfAdmission   int     `json:"year_of_admission"`
	IsVerified        bool    `json:"is_verified"`
	VerificationToken *string `json:"verification_token,omitempty"`
	StudentType       string  `json:"student_type"`
	RefreshToken      string  `json:"refresh_token,omitempty"`
	Role              string  `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	RollNum         string `json:"rollnum" binding:"required"`
	YearOfAdmission int    `json:"year_of_admission" binding:"required"`
	Branch          string `json:"branch"`
	StudentType     string `json:"student_type"`
}

type LoginResponse struct {
	Email string `json:"email"`
}
