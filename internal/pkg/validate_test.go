package pkg

import (
	"fmt"
	"testing"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
)

func TestValidateRegisterData(t *testing.T) {
	user := models.RegisterRequest{
		RollNum: "CO21314",
		Email: "co21314@ccet.ac.in",
		Name: "Charan",
		Password: "123456",
		YearOfAdmission: 2021,
		Branch: "Computer Science and Engineering",
		StudentType: "Regular",
	}

	result := ValidateRegisterData(user)

	if result != nil {
		t.Errorf("Expected no error, but got %v", result)
	}
}

func TestValidateRegisterDataInvalid(t *testing.T) {
	user := models.RegisterRequest{
		RollNum: "CO21314",
		Email: "co21545@ccet.ac.in",
		Name: "Charan",
		Password: "123456",
		YearOfAdmission: 2021,
		Branch: "Computer Science and Engineering",
		StudentType: "Regular",
	}

	result := ValidateRegisterData(user)

	if result.Error() != fmt.Errorf("Invalid Email").Error() {
		t.Errorf(result.Error())
	}
}
