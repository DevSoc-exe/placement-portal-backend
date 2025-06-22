package pkg

import (
	"testing"
)



func TestOTPEmail(t *testing.T) {

	email := CreateOTPEmail(1234, "charan", "co21314@ccet.ac.in")

	err := email.SendEmail()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err.Error())
	}
}
