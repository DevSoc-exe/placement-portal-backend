package dto

type StudentApplicationDTO struct {
	Id            string `json:"id"`
	ApplicationId string `json:"application_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Branch        string `json:"branch"`
	RollNum       string `json:"rollnum"`
	Gender        string `json:"gender"`
	Role          string `json:"role"`
	RoleId        string `json:"role_id"`
	AppliedAt     string `json:"applied_at"`
	IsPlaced      int    `json:"is_placed"`
}
