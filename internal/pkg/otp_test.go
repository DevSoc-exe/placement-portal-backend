package pkg

import (
	"log"
	"testing"
)

func TestOTPGeneration(t *testing.T) {

	otp, err := createOTP(NUMBER_OF_DIGITS)
	if err != nil {
		t.Error("Error creating OTP", err.Error())
	}

	log.Println(otp)
	result := otp/1000000 == 0

	if !result {
		t.Errorf("Expected no error, but got %v, %d", result, otp)
	}
}

// func TestOTPToken(t *testing.T) {
// 	otp, err := CreateOTP(NUMBER_OF_DIGITS)
// 	if err != nil {
// 		t.Errorf("Error creating OTP", err.Error())
// 	}

// 	_, err := CreateOtpJwt(otp)
// 	if err != nil {
// 		t.Errorf("Error Creating Token")
// 	}
// }
