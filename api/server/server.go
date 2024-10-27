package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DevSoc-exe/placement-portal-backend/api/handlers"
	"github.com/DevSoc-exe/placement-portal-backend/api/middleware"
	"github.com/DevSoc-exe/placement-portal-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Srv *gin.Engine
	Str models.Store
}

func CreateServer(server *gin.Engine, str models.Store) *Server {
	return &Server{
		Srv: server,
		Str: str,
	}
}

type Database struct {
	DB *sql.DB
}

func (server *Server) StartServer() {
	server.Srv.Use(middleware.CORSMiddleware)
	AddRoutes(server)

	if err := server.Srv.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	fmt.Println("Started server on port 8080.")
}

func AddRoutes(s *Server) {
	server := s.Srv.Group("/api")

	server.GET("/health", HealthCheck)

	//* Public User APIs
	server.PUT("/user/verify/:uid", handlers.HandleUserVerification(s.Str))

	//* Auth APIs
	server.GET("/refresh", middleware.RefreshToken(), handlers.HandleRefreshToken(s.Str))
	server.POST("/login", handlers.Login(s.Str))
	server.GET("/otp", handlers.HandleGetOTP(s.Str))
	server.POST("/signup", handlers.Register(s.Str))
	server.POST("/logout", handlers.HandleLogoutUser(s.Str))

	//* Job Posting API
	server.GET("/jobs/getDrive", handlers.HandleGetDriveUsingID(s.Str))
	server.DELETE("/jobs/delDrive", handlers.HandleDeleteDrive(s.Str))
	server.POST("/jobs/addNewDrive", handlers.HandleCreateNewDrive(s.Str))

	protectedServer := server.Group("/")
	protectedServer.Use(middleware.AuthMiddleware())
	{
		protectedServer.GET("/user", handlers.HandleGetUserdata(s.Str))
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
