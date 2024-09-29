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
}
