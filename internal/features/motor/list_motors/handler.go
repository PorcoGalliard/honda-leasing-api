package list_motors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(c *gin.Context) {
	var req ListMotorsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters: "+err.Error())
		return
	}

	if err := h.validateRequest(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.ListMotors(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch motors: "+err.Error())
		return
	}

	message := "Motors retrieved successfully"
	if req.GroupByType {
		message = "Motors grouped by type retrieved successfully"
	}

	response.Success(c, http.StatusOK, message, result)
}

func (h *Handler) validateRequest(req ListMotorsRequest) error {
	if req.MotorType != "" {
		validTypes := map[string]bool{
			"Sport":   true,
			"Matic":   true,
			"Classic": true,
			"Maxi":    true,
			"Bebek":   true,
		}
		if !validTypes[req.MotorType] {
			return response.NewValidationError("motor_type must be one of: Sport, Matic, Classic, Maxi, Bebek")
		}
	}

	if req.Status != "" {
		validStatuses := map[string]bool{
			"ready":    true,
			"booked":   true,
			"leased":   true,
			"returned": true,
			"repo":     true,
		}
		if !validStatuses[req.Status] {
			return response.NewValidationError("status must be one of: ready, booked, leased, returned, repo")
		}
	}

	if req.MinPrice < 0 {
		return response.NewValidationError("min_price must be >= 0")
	}

	if req.MaxPrice < 0 {
		return response.NewValidationError("max_price must be >= 0")
	}

	if req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		return response.NewValidationError("min_price must be less than or equal to max_price")
	}

	return nil
}
