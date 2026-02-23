package credit_simulation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(c *gin.Context) {
	var req CreditSimulationRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	// Call service
	result, err := h.service.Simulate(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMotorNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrDPExceedsPrice):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Gagal memproses simulasi kredit: "+err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Simulasi kredit berhasil dihitung", result)
}
