package models

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

	DeleteJobUsingDriveID(driveID string) error
	GetJobPostingUsingDriveID(driveID string) (interface{}, error)
	CreateNewDriveUsingObject(driveData DriveBody) error

	GetRolesUsingDriveID(driveID string) ([]Role, error)

	//* Company
	AddNewCompany(company *Company) error
	GetAllCompanies(args ...string) ([]Company, error)
	GetAllCompaniesForUser(args ...string) ([]CompanyResponse, error)

	AddStudentData(user *StudentData) error
	GetStudentDataByID(id string) (*StudentData, error)
	UpdateStudentData(user *StudentData) error
	DeleteStudentData(id string) error
	GetAllStudentData(args ...string) ([]*StudentData, error)
}
