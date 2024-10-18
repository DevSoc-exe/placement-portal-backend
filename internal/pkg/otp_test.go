package pkg

import (
	"log"
	"testing"
)

func TestOTPGeneration(t *testing.T) {

	otp := CreateOTP()
	log.Println(otp)
	result := otp/1000000 == 0

	if !result {
		t.Errorf("Expected no error, but got %v", result)
	}
}

func TestOTPToken(t *testing.T) {
	otp := CreateOTP()
	_, err := CreateOTPToken(otp)
	if err != nil {
		t.Errorf("Error Creating Token")
	}
}
