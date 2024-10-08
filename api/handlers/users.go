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
