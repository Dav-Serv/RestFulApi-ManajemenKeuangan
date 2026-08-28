package utils

import "github.com/gin-gonic/gin"

// success mengirim response suskses dengan format seragam.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"success": 	true,
		"message": 	message,
		"data":		data,
	})
}

// Error mengirim response error dengan format seragam.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}