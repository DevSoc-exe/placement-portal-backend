package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-CSRF-Token, Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		tokenString, err := c.Cookie("auth_token")
		if err != nil || tokenString == "" {
			fmt.Println("NO AUTH COOKIE", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		token, err := pkg.ValidateJWT(tokenString)
		if err != nil {
			if err.Error() == "Token is expired" {
				c.JSON(http.StatusNotAcceptable, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}

			if err == jwt.ErrSignatureInvalid {
				fmt.Println("JWT VALIDATION ERR")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("CANNOT EXTRACT TOKEN")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["userID"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract refresh token from request header
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Refresh token missing"})
			c.Abort()
			return
		}

		// Validate refresh token
		token, err := pkg.ValidateJWT(refreshToken)
		if err != nil {
			if err.Error() == "Token is expired" {
				c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
				c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			c.Abort()
			return
		}

		// Set user information in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			c.Abort()
			return
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		currentTime := time.Now().Local()

		if currentTime.After(expirationTime) {
			// return true
			c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
			c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("refresh_token", refreshToken)
		c.Set("userID", claims["userID"])
		c.Set("issuedAt", claims["issuedAt"])
		c.Next()
	}
}
