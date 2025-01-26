package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/DevSoc-exe/placement-portal-backend/internal/pkg"
	"github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
)

func HandleGetOTP(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		email := c.Query("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// auth, err := h.store.GetUserByEmail(req.Email)
		user, err := s.GetUserByEmail(email)
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err.Error()})
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not Verified"})
			return
		}

		otp, otpString, err := pkg.CreateOtpJwt()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = s.SaveOTP(otpString, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Save User Token in Database.", "message": err.Error()})
			return
		}

		mail := pkg.CreateOTPEmail(otp, user.Name, user.Email)
		err = mail.SendEmail()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Send OTP Email.", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func Login(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.OTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user, err := s.GetUserByEmail(req.Email)
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err.Error()})
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not Verified"})
			return
		}

		token := user.Otp
		otp, err := pkg.CheckOTPToken(token.String)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if otp != req.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		// if err := bcrypt.CompareHashAndPassword([]byte(user.Otp), []byte(string(req.OTP))); err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		// 	return
		// }

		err = s.ClearOTP(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear OTP"})
		}

		// otp, err := pkg.CheckOTPToken(token)

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token",
				"message": err.Error()})
			return
		}

		c.Set("csrf_token", csrfToken)
		c.Header("X-Csrf-Token", csrfToken)

		domain := os.Getenv("DOMAIN")
		secure := true
		if domain == "" {
			secure = false
			domain = "localhost"
		} else {
			domain = ".classikh.me"
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

		//! No need to create a password due to OTP Login
		// hash, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		csrfToken, err := pkg.GenerateCSRFToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}
		verification_token := sql.NullString{
			String: csrfToken,
		}

		userId := nanoid.New()
		newauth := models.User{
			ID:                userId,
			Email:             auth.Email,
			Name:              auth.Name,
			PhoneNumber:       auth.PhoneNumber,
			Gender:            auth.Gender,
			YearOfAdmission:   auth.YearOfAdmission,
			VerificationToken: verification_token,
			RollNumber:        auth.RollNum,
			Branch:            auth.Branch,
			StudentType:       auth.StudentType,
			Role:              "STUDENT",
		}

		if err := s.CreateUser(&newauth); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mail := pkg.CreateMailMessageWithVerificationToken(csrfToken, userId, newauth.Email)
		err = mail.SendEmail()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func HandleGetUserdata(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in context"})
			return
		}
		fmt.Println(id)

		user, err := s.GetUserByID(id.(string))
		userRes := &models.UserResponse{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			Gender:          user.Gender,
			RollNumber:      user.RollNumber,
			YearOfAdmission: user.YearOfAdmission,
			Branch:          user.Branch,
			StudentType:     user.StudentType,
			IsOnboarded:     user.IsOnboarded,
		}
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": userRes, "success": true})
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
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Verify User Token!",
				"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func HandleLogoutUser(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found."})
			return
		}
		fmt.Println(id)

		err := s.RevokeUserRefreshToken(id.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occured"})
			return
		}

		domain := os.Getenv("DOMAIN")
		secure := true
		if domain == "" {
			secure = false
			domain = "localhost"
		} else {
			domain = ".classikh.me"
		}
		c.SetCookie("auth_token", "", -1, "/", domain, secure, true)
		c.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
		c.JSON(http.StatusNoContent, nil)
	}
}

func HandleRefreshToken(s models.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve auth by ID
		refreshToken, _ := c.Get("refresh_token")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token not found in context"})
			return
		}

		id, _ := c.Get("userID")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
			return
		}

		auth, err := s.GetUserByID(id.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if auth.RefreshToken != refreshToken.(string) {
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
		} else {
			domain = ".classikh.me"
		}

		c.SetCookie("auth_token", accessToken, 3600*24, "/", domain, secure, true)

		c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
	}
}
