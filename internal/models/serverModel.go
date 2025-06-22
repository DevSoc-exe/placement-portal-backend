package models

import (
	"context"
	"database/sql"
	"github.com/DevSoc-exe/placement-portal-backend/internal/models/dto"
)

type Store interface {
	CreateUser(user *User) error
	VerifyUser(userId string, token string) error
	UpdateUserRefreshToken(refresh_token, userId string) error
	RevokeUserRefreshToken(refresh_token string) error
	SaveOTP(otpToken string, userId string) error
	ClearOTP(userID string) error

	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetAllStudents(args ...string) ([]*UserResponse, error)
	ToggleUserRole(id string) error
	GetAllDrivesForUser() ([]DriveResponse, error)
	GetUserMailsByBranchesAboveCGPA(branches []string, cgpaLimit float32) ([]string, error)

	//* Applications
	ApplyForDrive(userID, roleID, driveID string) error
	GetAppliedRole(userID string, driveId string) (*Role, error)
	GetDriveApplicantsForRole(roleID, required_data, driveID string) (*sql.Rows, []string, error)

	DeleteJobUsingDriveID(driveID string) error
	GetJobPostingUsingDriveID(driveID string) (*Drive, error)
	CreateNewDriveUsingObject(driveData Drive) (string, error)

	GetRolesUsingDriveID(driveID string) ([]Role, error)
	GetApplicantsForDrive(driveID string) ([]dto.StudentApplicationDTO, error)

	//* Company
	AddNewCompany(company *Company) error
	GetAllCompanies(args ...string) ([]Company, error)
	GetAllCompaniesForUser(args ...string) ([]CompanyResponse, error)
	GetCompanyUsingCompanyID(companyID string) (*CompanyResponse, error)

	AddStudentData(user *StudentData) error
	GetStudentDataByID(id string) (*StudentData, error)
	UpdateStudentData(context context.Context, user *StudentData) error
	DeleteStudentData(id string) error
	GetAllStudentData(args ...string) ([]*StudentData, error)

	//* Applications Table
	MarkStudentAsPlaced(updatedById string, applicationId string) error
}
