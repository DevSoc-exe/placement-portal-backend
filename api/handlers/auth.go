package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	"github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	RollNum         string `json:"rollnum" binding:"required"`
	YearOfAdmission int    `json:"year_of_admission" binding:"required"`
}

type LoginResponse struct {
	Email string `json:"email"`
}

func Login(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// auth, err := h.store.GetUserByEmail(req.Email)
		user, err := s.GetUserByEmail(req.Email)
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err})
			return
		}
		log.Println(user)

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or password"})
			return
		}

		csrfToken, err := pkg.GenerateCSRFToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}

		// Generate tokens
		accessToken, err := pkg.CreateAccessToken(user, csrfToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create access token"})
			return
		}

		refreshToken, err := pkg.CreateRefreshToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create refresh token"})
			return
		}

		// // Update auth with the new refresh token
		// auth.RefreshToken = refreshToken
		// if err := h.store.UpdateUserRefreshToken(refreshToken, auth.ID); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
		// 	return
		// }

		c.Set("csrf_token", csrfToken)
		c.Header("X-Csrf-Token", csrfToken)

		domain := os.Getenv("DOMAIN")
		secure := true
		if domain == "" {
			secure = false
			domain = "localhost"
		}

		c.SetCookie("auth_token", accessToken, 3600*24, "/", domain, secure, true)
		c.SetCookie("refresh_token", refreshToken, 3600*24, "/", domain, secure, true)

		resp := LoginResponse{
			Email: user.Email,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func Register(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth RegisterRequest
		if err := c.ShouldBindJSON(&auth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newauth := models.User{
			ID:              nanoid.New(),
			Email:           auth.Email,
			Password:        string(hash[:]),
			Name:            auth.Name,
			YearOfAdmission: auth.YearOfAdmission,
			RollNumber:      auth.RollNum,
		}

		if err := s.CreateUser(&newauth); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func HandleLogoutUser(s models.Store) gin.HandlerFunc {
	return func (c *gin.Context) {
		authID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
			return
		}

		fmt.Println(authID)

		// var id int
		// if floatID, ok := authID.(float64); ok {
		// 	id = int(floatID)
		// } else {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "auth_id is not a valid number"})
		// 	return
		// }

		// err := h.store.RevokeUserRefreshToken(id)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occured"})
		// 	return
		// }

		domain := os.Getenv("DOMAIN")
		secure := true
		if domain == "" {
			secure = false
			domain = "localhost"
		}

		c.SetCookie("auth_token", "", -1, "/", domain, secure, true)
		c.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
	}
}
