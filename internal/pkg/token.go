package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/config"
	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/dgrijalva/jwt-go"
)

const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

func CreateAccessToken(user *models.User, csrf string) (string, error) {
	claims := jwt.MapClaims{
		"exp":      time.Now().Local().Add(time.Minute * 15).Unix(), // Token expiry time
		"userID":   user.ID,
		"email":    user.Email,
		"issuedAt": time.Now().Local().Unix(),
		"role":     user.Role,
		"csrf":     csrf,
	}

	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	authTokenString, err := authJwt.SignedString(config.SignKey)
	if err != nil {
		return "", err
	}
	return authTokenString, nil
}

func CreateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"exp":      time.Now().Local().Add(time.Hour * 24).Unix(), // Set a longer expiration time
		"userID":   user.ID,
		"email":    user.Email,
		"role":     user.Role,
		"issuedAt": time.Now().Local().Unix(),
	}

	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token using the RSA private key
	refreshTokenString, err := token.SignedString(config.SignKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for signing
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return config.VerifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

func IsTokenExpired(tokenString string) bool {
	token, err := ValidateJWT(tokenString)
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	currentTime := time.Now()

	if currentTime.After(expirationTime) {
		return true
	}

	return false
}

func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
