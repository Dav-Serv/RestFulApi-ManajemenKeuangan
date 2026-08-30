package routes

import (
	"net/http"

	"RestFulApi-ManajemenKeuangan/internal/handlers"
	"RestFulApi-ManajemenKeuangan/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ---- Auth (public) ----
	auth := r.Group("/auth")
	{
		auth.GET("/google/login", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.Googlecallback)
	}

	// ---- API (protected pakai JWT) ----
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/auth/me", handlers.Me)

		trx := api.Group("/transactions")
		{
			trx.POST("", handlers.CreateTransaction)
			trx.GET("", handlers.GetTransactions)
			trx.GET("/summary", handlers.GetSummary)
			trx.GET("/:id", handlers.GetTransactionById)
			trx.PUT("/:id", handlers.UpdateTransaction)
			trx.DELETE("/:id", handlers.DeleteTransaction)
		}
	}

	return r
}