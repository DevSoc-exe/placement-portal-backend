package pkg

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//* Number of digits in the OTP
const NUMBER_OF_DIGITS = 6

func CreateOtpJwt() (int, string, error) {
	otp, err := createOTP(NUMBER_OF_DIGITS)
	if err != nil {
		return 0, "", err
	}

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
			return 0, fmt.Errorf("token expired")
		}

		if err == jwt.ErrSignatureInvalid {
			return 0, fmt.Errorf("invalid or expired token")
		}

		return 0, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
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

func createOTP(numberOfDigits int) (int, error) {
    maxLimit := int64(int(math.Pow10(numberOfDigits)) - 1)
    lowLimit := int(math.Pow10(numberOfDigits - 1))

    randomNumber, err := rand.Int(rand.Reader, big.NewInt(maxLimit))
    if err != nil {
        return 0, err
    }
    randomNumberInt := int(randomNumber.Int64())

    // Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
    if randomNumberInt <= lowLimit {
        randomNumberInt += lowLimit
    }

    // Never likely to occur, kust for safe side.
    if randomNumberInt > int(maxLimit) {
        randomNumberInt = int(maxLimit)
    }
    return randomNumberInt, nil
}
