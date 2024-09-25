package models

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password		string `json:"password"`
	RollNumber      string `json:"rollnum"`
	YearOfAdmission int    `json:"year_of_admission"`
}
