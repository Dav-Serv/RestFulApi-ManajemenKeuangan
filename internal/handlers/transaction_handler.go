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
		utils.Error(c, http.StatusInternalServerError, "gagal menyimpan transaksi")
		return
	}

	utils.Success(c, http.StatusCreated, "transaksi berhasil dibuat", trx)
}

// GetTransactions - GET /api/transactions
// Mendukung query param: type, category, start_date, end_date, page, limit
func GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")

	query := database.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)

	if t := c.Query("type"); t != "" {
		query = query.Where("type = ?", t)
	}
	if cat := c.Query("category"); cat != "" {
		query = query.Where("category = ?", cat)
	}
	if start := c.Query("start_date"); start != "" {
		query = query.Where("date >= ?", start)
	}
	if end := c.Query("end_date"); end != "" {
		query = query.Where("date <= ?", end)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	var transaction []models.Transaction
	if err := query.Order("date desc").Offset(offset).Limit(limit).Find(&transaction).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil data transaksi", gin.H{
		"transactions": transaction,
		"pagination": gin.H{
			"page":	 page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetTransactionByID - GET /api/transactions/:id
func GetTransactionById(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var trx models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&trx).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "transaksi tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil transaksi", trx)
}

// UpdateTransaction - PUT /api/transactions/:id
func UpdateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var trx models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&trx).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "transaksi tidak ditemukan")
		return
	}

	var input models.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid" + err.Error())
		return
	}

	trx.Type = input.Type
	trx.Category = input.Category
	trx.Amount = input.Amount
	trx.Description = input.Description
	if input.Date != nil {
		trx.Date = *input.Date
	}

	if err := database.DB.Save(&trx).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal memperbarui transaksi")
		return
	}

	utils.Success(c, http.StatusOK, "transaksi berhasil diperbarui", trx)
}

// DeleteTransaction - DELETE /api/transactions/:id
func DeleteTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var trx models.Transaction
	if err := database.DB.Where("id = ? and user_id = ?", id, userID).First(&trx).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "transaksi tidak ditemukan")
		return
	}

	if err := database.DB.Delete(&trx).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal menghapus transaksi")
		return
	}

	utils.Success(c, http.StatusOK, "transaksi berhasil dihapus", nil)
}

// GetSummary - GET /api/transactions/summary
// Ringkasan total pemasukan, pengeluaran, dan saldo milik user yang login.
func GetSummary(c *gin.Context) {
	userID := c.GetUint("user_id")

	var totalIncome, totalExpense float64

	database.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, models.TypeIncome).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)

	database.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, models.TypeExpense).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	summary := models.SummaryResponse{
		TotalIncome:	totalIncome,
		TotalExpense:	totalExpense,
		Balance: 		totalIncome - totalExpense,
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil ringkasan", summary)
}