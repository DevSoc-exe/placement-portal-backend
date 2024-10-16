package pkg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)


func CreateOtpJwt() (int, string, error) {
	otp := createOTP()

	claims := jwt.MapClaims{
		"exp":      time.Now().Local().Add(time.Minute * 15).Unix(), // Token expiry time
		"issuedAt": time.Now().Local().Unix(),
		"otp":      otp,
	}

	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	authTokenString, err := authJwt.SignedString(config.SignKey)
	if err != nil {
		return 0, "", err
	}
	return otp, authTokenString, nil
}

func CreateOTPHash(otp int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(string(otp)), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckOTPToken(tokenString string) (int, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		if err.Error() == "Token is expired" {
			return 0, fmt.Errorf("Token expired")
		}

		if err == jwt.ErrSignatureInvalid {
			return 0, fmt.Errorf("Invalid or expired token")
		}

		return 0, fmt.Errorf("Invalid Token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("Invalid Token")
	}

	otpValue, ok := claims["otp"]
	if !ok {
		return 0, fmt.Errorf("OTP not found in token")
	}

	switch otp := otpValue.(type) {
	case float64:
		return int(otp), nil
	case int:
		return otp, nil
	default:
		return 0, fmt.Errorf("OTP is not a valid number")
	}
}

func createOTP() int {
	place := 1
	var result = 0

	for i := 1; i <= 6; i++ {
		num := place * rand.Intn(10)
		result = result + num
		place = place * 10
	}

	return result
}
