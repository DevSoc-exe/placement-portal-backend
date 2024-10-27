package responses

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// func CreateApiResponse(success bool, message string, data struct{}) *ApiResponse {
// 	return &ApiResponse{
// 		Success: success,
// 		Message: message,
// 		Data:    data,
// 	}
// }

func (api *ApiResponse) MapApiResponse(c *gin.Context, statusCode int) {
	c.JSON(statusCode, gin.H{
		"Success": api.Success,
		"Message": api.Message,
		"Data":    api.Data,
	})
}
