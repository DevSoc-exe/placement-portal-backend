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
	server.GET("/admin/user", handlers.HandleGetAllStudents(s.Str))
	server.GET("/company", handlers.HandleGetAllCompanies(s.Str))
	server.POST("/company", handlers.HandleCreateNewCompany(s.Str))
	server.GET("/getCompanyFromID", handlers.HandleGetCompanyFromID(s.Str))

	//* Auth APIs
	server.GET("/refresh", middleware.RefreshToken(), handlers.HandleRefreshToken(s.Str))
	server.POST("/login", handlers.Login(s.Str))
	server.GET("/otp", handlers.HandleGetOTP(s.Str))
	server.POST("/signup", handlers.Register(s.Str))
	server.POST("/logout", handlers.HandleLogoutUser(s.Str))

	// *temporarily added as public routes --------------------
	server.POST("/jobs/addNewDrive", handlers.HandleCreateNewDrive(s.Str))
	server.POST("/jobs/drive/applicant", handlers.HandleGetDriveApplicantsForRole(s.Str))
	server.PUT("/admin/user/role/:id", handlers.HandleToggleUserRole(s.Str))
	server.GET("/jobs/drive/applied/:id", handlers.HandleGetApplicantsForDrive(s.Str))


	//* Job Posting API
	//! TO BE REMOVED TO PROTECTED API
	server.GET("/jobs/getDrive", handlers.HandleGetDriveUsingID(s.Str))

	//* Admin APIs
	adminServer := server.Group("/")
	adminServer.Use(middleware.AuthMiddleware(), middleware.CheckAdmin())
	{
		// adminServer.PUT("/admin/user/role/:id", handlers.HandleToggleUserRole(s.Str))

		adminServer.DELETE("/jobs/delDrive", handlers.HandleDeleteDrive(s.Str))
		// adminServer.POST("/jobs/addNewDrive", handlers.HandleCreateNewDrive(s.Str))

		//* Student Data APIs for admin
		adminServer.GET("/admin/user/data", handlers.HandleGetAllStudentData(s.Str))
		adminServer.GET("/admin/user/data/:id", handlers.HandleGetStudentDataByID(s.Str))


	}

	//* User APIs
	userServer := server.Group("/")
	userServer.Use(middleware.AuthMiddleware())
	{
		//* Application routes
		//! TO BE MOVED TO ADMIN
		userServer.POST("/drive/placed/:application_id", handlers.HandleMarkStudentAsPlaced(s.Str))


		userServer.GET("/user", handlers.HandleGetUserdata(s.Str))

		//* Drive APIs for user
		userServer.GET("/user/drive", handlers.HandleGetDrivesForUser(s.Str))
		userServer.GET("/user/drive/:id", handlers.HandleGetDriveUsingID(s.Str))
		userServer.POST("/user/drive", handlers.HandleApplyToDrive(s.Str))

		//* Company APIs for user
		userServer.GET("/user/company", handlers.HandleGetCompaniesForUser(s.Str))

		//* Student Data APIs for user
		userServer.POST("/user/data", handlers.HandleAddNewStudentData(s.Str))
		userServer.GET("/user/data", handlers.HandleGetStudentData(s.Str))
		userServer.DELETE("/user/data", handlers.HandleDeleteStudentData(s.Str))
		userServer.PUT("/user/data", handlers.HandleUpdateStudentData(s.Str))
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
