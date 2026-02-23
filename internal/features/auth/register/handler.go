package register

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedErrors "github.com/nanasuryana335/honda-leasing-api/internal/shared/errors"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, sharedErrors.ErrDuplicateEntry):
			response.Error(c, http.StatusConflict, "Nomor HP sudah terdaftar")
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Registrasi berhasil", result)
}
