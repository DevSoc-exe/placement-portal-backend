package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
)

func ValidateRegisterData(data models.RegisterRequest) error {
	rollNum := data.RollNum

	matches, err := regexp.MatchString(`^(LCO|MCO|CO)\d{5}$`, rollNum)
	if err != nil {
		return err
	}

	if !matches {
		return fmt.Errorf("Invalid RollNumber")
	}

	studentType := string(rollNum[0]) + string(rollNum[1])
	validStudentTypes := map[string]string{
		"CO":  "Regular",
		"LCO": "LEET",
		"MCO": "PU MEET",
	}

	expectedType, ok := validStudentTypes[studentType]
	if !ok {
		return fmt.Errorf("Invalid Student Type")
	}

	if expectedType != data.StudentType {
		return fmt.Errorf("Student Type Mismatch")
	}

	gender := data.Gender
	if gender != "MALE" && gender != "FEMALE" && gender != "OTHERS" {
		return fmt.Errorf("Invalid Gender")
	}

	year := 2000 + (int(rollNum[2]-'0') * 10) + int(rollNum[3]-'0')
	if year != data.YearOfAdmission {
		return fmt.Errorf("Year of Admission Mismatch")
	}

	branchCode := string(rollNum[4])
	validBranchCodes := map[string]string{
		"3": "Computer Science and Engineering",
		"5": "Electronics and Communications Engineering",
		"1": "Mechanical Engineering",
		"2": "Civil Engineering",
	}

	expectedBranch, ok := validBranchCodes[branchCode]
	if !ok {
		return fmt.Errorf("Invalid Branch")
	}

	if expectedBranch != data.Branch {
		return fmt.Errorf("Student Branch Mismatch")
	}

	email := strings.ToLower(rollNum) + "@ccet.ac.in"
	if email != data.Email {
		return fmt.Errorf("Invalid Email")
	}

	return nil
}
