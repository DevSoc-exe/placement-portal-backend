package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DevSoc-exe/placement-portal-backend/api/server"
	"github.com/DevSoc-exe/placement-portal-backend/internal/config"
	"github.com/DevSoc-exe/placement-portal-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func init() {
	config.InitEnv()

	// No need to create keys every time during development
	env := os.Getenv("ENVIRONMENT")
	if env == "PRODUCTION" {
		config.CreateKeys()
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		panic("database connection string not found")
	}
	gin.SetMode(mode)

	config.InitJWT()
}

func main() {
	r := gin.Default()
	fmt.Println("in to server")
	dsn := os.Getenv("DB_CONN")
	if dsn == "" {
		panic("database connection string not found")
	}

	db, err := database.ConnectToDB(dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	log.Println("Connected to the database!")
	defer db.Close()

	store := database.CreateDatabase(db)
	err = store.InitDB()
	if err != nil {
		log.Fatalln(err)
	}

	srv := server.CreateServer(r, store)
	srv.StartServer()
}
