package pkg

import (
	"math/rand"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/config"
	"github.com/dgrijalva/jwt-go"
)

func CreateOTP() int {
	place := 1
	var result = 0

	for i := 1; i <= 6; i++ {
		num := place * rand.Intn(10)
		result = result + num
		place = place * 10
	}

	return result
}

func CreateOTPToken(otp int) (string, error) {
	claims := jwt.MapClaims{
		"exp":      time.Now().Local().Add(time.Minute * 5).Unix(),
		"otp":      otp,
		"issuedAt": time.Now().Local().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	refreshTokenString, err := token.SignedString(config.SignKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}
