package create_order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/middleware"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Tidak dapat mengidentifikasi user")
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	result, err := h.service.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMotorNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrMotorNotAvailable):
			response.Error(c, http.StatusConflict, err.Error())
		case errors.Is(err, ErrDPExceedsPrice):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrNIKRequired):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Gagal membuat order: "+err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Order berhasil dibuat", result)
}
