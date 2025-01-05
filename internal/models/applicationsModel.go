package models

type Application struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id" binding:"required"`
	DriveID   string `json:"drive_id" binding:"required"`
	AppliedAt string `json:"applied_at"`
}
