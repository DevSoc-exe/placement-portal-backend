package models

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
	ID               string  `json:"id"`
	CompanyID        string  `json:"company_id"`
	Company          Company `json:"company" gorm:"foreignKey:CompanyID"`
	DateOfDrive      string  `json:"drive_date"`
	DriveDuration    int     `json:"drive_duration"`
	Roles            []Role  `json:"roles" gorm:"foreignKey:DriveID"`
	Location         string  `json:"location"`
	Qualifications   string  `json:"qualifications"`
	PointsToNote     string  `json:"points_to_note"`
	JobDescription   string  `json:"job_description"`
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
	CompanyID      string `json:"company_id"`
	DateOfDrive    string `json:"drive_date"`
	DriveDuration  int    `json:"drive_duration"`
	Roles          []Role `json:"roles"`
	Location       string `json:"location"`
	Qualifications string `json:"qualifications"`
	PointsToNote   string `json:"points_to_note"`
	JobDescription string `json:"job_description"`
}
