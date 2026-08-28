package handlers

import (
	"net/http"
	"strconv"
	"time"

	"RestFulApi-ManajemenKeuangan/internal/database"
	"RestFulApi-ManajemenKeuangan/internal/models"
	"RestFulApi-ManajemenKeuangan/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateTransaction - POST /api/transactions
func CreateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input models.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid:" + err.Error())
		return
	}

	date := time.Now()
	if input.Date != nil {
		date = *input.Date
	}

	trx := models.Transaction{
		UserID: 		userID,
		Type: 			input.Type,
		Category: 		input.Category,
		Amount: 		input.Amount,
		Description: 	input.Description,
		Date: 			date,
	}

	if err := database.DB.Create(&trx).Error; err != nil {
		utils.Error()
	}
}