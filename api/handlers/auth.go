package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	"github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
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

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not Verified"})
			return
		}

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

		if err := s.UpdateUserRefreshToken(refreshToken, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
			return
		}

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

		resp := models.LoginResponse{
			Email: user.Email,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func Register(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth models.RegisterRequest
		if err := c.ShouldBindJSON(&auth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate data
		err := pkg.ValidateRegisterData(auth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		csrfToken, err := pkg.GenerateCSRFToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}

		userId := nanoid.New()
		newauth := models.User{
			ID:                userId,
			Email:             auth.Email,
			Password:          string(hash[:]),
			Name:              auth.Name,
			YearOfAdmission:   auth.YearOfAdmission,
			RollNumber:        auth.RollNum,
			Branch:            auth.Branch,
			StudentType:       auth.StudentType,
			VerificationToken: &csrfToken,
			Role:              "STUDENT",
		}

		if err := s.CreateUser(&newauth); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		emailBody := pkg.CreateMailMessageWithVerificationToken(csrfToken, userId)
		pkg.SendVerificationEmail(auth.Email, emailBody)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func HandleGetUserdata(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in context"})
			return
		}
		fmt.Println(email)

		user, err := s.GetUserByEmail(email.(string))
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": user, "success": true})
	}
}

type VerifyRequestBody struct {
	Token string `json:"token" binding:"required"`
}

func HandleUserVerification(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")
		if uid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
			return
		}

		var token VerifyRequestBody
		if err := c.ShouldBindJSON(&token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := s.VerifyUser(uid, token.Token)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Verify User Token!",
				"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func HandleLogoutUser(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("userID")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
			return
		}

		err := s.RevokeUserRefreshToken(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occured"})
			return
		}

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

func HandleRefreshToken(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve auth by ID
		refreshToken := c.GetString("refresh_token")
		id := c.GetString("userID")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
			return
		}

		auth, err := s.GetUserByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if auth.RefreshToken != refreshToken {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token did not match"})
			return
		}

		// Generate a new access token
		accessToken, err := pkg.CreateAccessToken(auth, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create access token"})
			return
		}

		domain := os.Getenv("DOMAIN")
		secure := true
		if domain == "" {
			secure = false
			domain = "localhost"
		}

		c.SetCookie("auth_token", accessToken, 3600*24, "/", domain, secure, true)

		c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
	}
}
