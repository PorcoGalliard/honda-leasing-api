package update_order_status

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/leasing"
	"github.com/nanasuryana335/honda-leasing-api/internal/models"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/middleware"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Handle(c *gin.Context) {

	roleName := getPrimaryRole(c)
	if roleName == "" {
		response.Error(c, http.StatusForbidden, "Tidak dapat menentukan role")
		return
	}

	contractIDStr := c.Param("contract_id")
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "contract_id tidak valid")
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	ctx := c.Request.Context()

	contract, err := findContractByID(ctx, h.db, contractID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data kontrak: "+err.Error())
		return
	}
	if contract == nil {
		response.Error(c, http.StatusNotFound, ErrContractNotFound.Error())
		return
	}

	if !isValidTransition(contract.Status, req.Status) {
		response.Error(c, http.StatusBadRequest,
			ErrInvalidTransition.Error()+": "+contract.Status+" → "+req.Status+" tidak diizinkan")
		return
	}

	if !isRoleAllowed(roleName, req.Status) {
		response.Error(c, http.StatusForbidden, ErrForbiddenTransition.Error())
		return
	}

	oldStatus := contract.Status
	lq := leasing.Use(h.db)
	lc := lq.LeasingContract
	_, err = lc.WithContext(ctx).
		Where(lc.ContractID.Eq(contractID)).
		Updates(map[string]interface{}{
			"status":     req.Status,
			"updated_at": time.Now(),
		})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal update status: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Status order berhasil diupdate", UpdateOrderStatusResponse{
		ContractID:     contract.ContractID,
		ContractNumber: contract.ContractNumber,
		StatusLama:     oldStatus,
		StatusBaru:     req.Status,
	})
}

func findContractByID(ctx context.Context, db *gorm.DB, contractID int64) (*models.LeasingContract, error) {
	lq := leasing.Use(db)
	lc := lq.LeasingContract
	contract, err := lc.WithContext(ctx).
		Where(lc.ContractID.Eq(contractID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return contract, nil
}

func isValidTransition(from, to string) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	return allowed[to]
}

func isRoleAllowed(role, targetStatus string) bool {
	allowed, ok := roleTransitionAllowed[role]
	if !ok {
		return false
	}
	return allowed[targetStatus]
}

func getPrimaryRole(c *gin.Context) string {
	roles, ok := middleware.GetRoles(c)
	if !ok {
		return ""
	}
	for _, r := range roles {
		if r != "CUSTOMER" {
			return r
		}
	}
	return ""
}
