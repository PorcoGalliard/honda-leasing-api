package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	basePath := viper.GetString("SERVER.BASE_PATH")
	api := router.Group(basePath)
	{
		// region routes endpoints
		users := api.Group("/user")
		{
			users.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "oke wir"})
			})
		}
	}
}
