package models

import "time"

type Company struct {
	CompanyID     string `json:"id"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	HRName        string `json:"hrName"`
	ContactEmail  string `json:"contactEmail"`
	ContactNumber string `json:"contactNumber"`
	LinkedIn      string `json:"linkedIn"`
	Website       string `json:"website"`
}

type CompanyResponse struct {
	CompanyID string `json:"id"`
	Name      string `json:"name"`
	Overview  string `json:"overview"`
	LinkedIn  string `json:"linkedIn"`
	Website   string `json:"website"`
}

type Drive struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	Company        Company   `json:"company" gorm:"foreignKey:CompanyID"`
	DateOfDrive    time.Time `json:"drive_date"`
	DriveDuration  int       `json:"drive_duration"`
	Roles          []Role    `json:"roles" gorm:"foreignKey:DriveID"`
	MinCGPA        float32   `json:"min_cgpa"`
	Deadline       time.Time `json:"deadline"`
	Location       string    `json:"location"`
	Qualifications string    `json:"qualifications"`
	PointsToNote   string    `json:"points_to_note"`
	JobDescription string    `json:"job_description"`
	DriveType      string    `json:"drive_type"`
	AppliedRole    Role      `json:"applied_role"`
	RequiredData   string    `json:"required_data"`
	Cse_allowed    bool      `json:"cse_allowed"`
	Ece_allowed    bool      `json:"ece_allowed"`
	Mech_allowed   bool      `json:"mech_allowed"`
	Civ_allowed    bool      `json:"civ_allowed"`
	Expired        bool      `json:"expired"`
}

type Role struct {
	ID          string `json:"id"`
	DriveID     string `json:"drive_id"`
	Title       string `json:"title"`
	StipendLow  int    `json:"stipend_low"`
	StipendHigh int    `json:"stipend_high"`
	SalaryLow   int    `json:"salary_low"`
	SalaryHigh  int    `json:"salary_high"`
}

type DriveBody struct {
	CompanyID       string  `json:"company_id"`
	DateOfDrive     time.Time  `json:"drive_date"`
	DriveDuration   int     `json:"drive_duration"`
	Roles           []Role  `json:"roles"`
	Deadline        time.Time  `json:"deadline"`
	Location        string  `json:"location"`
	Qualifications  string  `json:"qualifications"`
	PointsToNote    string  `json:"points_to_note"`
	JobDescription  string  `json:"job_description"`
	MinCGPA         float32 `json:"min_cgpa"`
	DriveType       string  `json:"drive_type"`
	RequiredData    string  `json:"required_data"`
	AllowedBranches string  `json:"allowed_branches"`
}

type DriveResponse struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	CompanyName    string    `json:"name"`
	DateOfDrive    time.Time `json:"drive_date"`
	DriveDuration  int       `json:"drive_duration"`
	Roles          []Role    `json:"roles"`
	Deadline       time.Time `json:"deadline"`
	Location       string    `json:"location"`
	Qualifications string    `json:"qualifications"`
	PointsToNote   string    `json:"points_to_note"`
	JobDescription string    `json:"job_description"`
	MinCGPA        float32   `json:"min_cgpa"`
	DriveType      string    `json:"drive_type"`
	AppliedRole    Role      `json:"applied_role"`
}
