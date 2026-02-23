package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/config"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/auth/login"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/auth/register"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/motor/credit_simulation"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/motor/list_motors"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/order/create_order"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/order/get_order_progress"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/staff/list_orders"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/staff/update_order_status"
	"github.com/nanasuryana335/honda-leasing-api/internal/features/staff/update_task_status"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/middleware"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
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

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			// Register - daftar customer baru
			registerHandler := register.NewHandler(register.NewService(register.NewRepository(db)))
			auth.POST("/register", registerHandler.Handle)

			// Login - masuk dengan nomor HP dan PIN
			loginHandler := login.NewHandler(login.NewService(login.NewRepository(db), cfg.JWT.Secret))
			auth.POST("/login", loginHandler.Handle)
		}

		// Motor routes (public)
		motors := api.Group("/motors")
		{
			// List motors - public endpoint for customers to browse
			listMotorsService := list_motors.NewService(list_motors.NewRepository(db))
			listMotorsHandler := list_motors.NewHandler(listMotorsService)
			motors.GET("", listMotorsHandler.Handle)

			// Credit simulation - public, untuk simulasi sebelum order
			creditSimService := credit_simulation.NewService(credit_simulation.NewRepository(db))
			creditSimHandler := credit_simulation.NewHandler(creditSimService)
			motors.POST("/credit-simulation", creditSimHandler.Handle)
		}

		// Protected routes - butuh JWT
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(cfg.JWT.Secret))
		{
			// Test JWT
			protected.GET("/tes-jwt", func(c *gin.Context) {
				userID, _ := middleware.GetUserID(c)
				roles, _ := middleware.GetRoles(c)
				c.JSON(http.StatusOK, gin.H{"user_id": userID, "roles": roles})
			})

			// Order routes
			orders := protected.Group("/orders")
			{
				// POST /orders - Submit order baru
				createOrderHandler := create_order.NewHandler(create_order.NewService(create_order.NewRepository(db)))
				orders.POST("", createOrderHandler.Handle)

				// GET /orders/:contract_id/progress - Lihat progress order
				getProgressHandler := get_order_progress.NewHandler(get_order_progress.NewRepository(db))
				orders.GET("/:contract_id/progress", getProgressHandler.Handle)
			}

			staff := protected.Group("/staff")
			staff.Use(middleware.RequireStaff())
			{
				listOrdersHandler := list_orders.NewHandler(list_orders.NewRepository(db))
				staff.GET("/orders", listOrdersHandler.Handle)

				// Update status order
				updateOrderStatusHandler := update_order_status.NewHandler(db)
				staff.PATCH("/orders/:contract_id/status", updateOrderStatusHandler.Handle)

				// Update status task
				updateTaskStatusHandler := update_task_status.NewHandler(db)
				staff.PATCH("/orders/:contract_id/tasks/:task_id", updateTaskStatusHandler.Handle)
			}
		}
	}
}
