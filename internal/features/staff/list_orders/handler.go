package list_orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/middleware"
	"github.com/nanasuryana335/honda-leasing-api/internal/shared/response"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}
func (h *Handler) Handle(c *gin.Context) {
	var req ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter tidak valid: "+err.Error())
		return
	}

	roleName := getPrimaryRole(c)
	if roleName == "" {
		response.Error(c, http.StatusForbidden, "Tidak dapat menentukan role")
		return
	}

	total, err := h.repo.CountOrders(c.Request.Context(), req, roleName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal menghitung total order: "+err.Error())
		return
	}

	rows, _, err := h.repo.FindOrders(c.Request.Context(), req, roleName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data order: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	orders := make([]OrderItem, 0, len(rows))
	for _, r := range rows {
		orders = append(orders, OrderItem{
			ContractID:      r.ContractID,
			ContractNumber:  r.ContractNumber,
			RequestDate:     r.RequestDate,
			Status:          r.Status,
			CustomerName:    r.CustomerName,
			CustomerPhone:   r.CustomerPhone,
			MotorMerk:       r.MotorMerk,
			MotorType:       r.MotorType,
			NilaiKendaraan:  r.NilaiKendaraan,
			DpDibayar:       r.DpDibayar,
			TenorBulan:      r.TenorBulan,
			CicilanPerBulan: r.CicilanPerBulan,
		})
	}

	response.Success(c, http.StatusOK, "Data order berhasil diambil", ListOrdersResponse{
		Orders: orders,
		Pagination: Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalRows:  total,
			TotalPages: calcTotalPages(total, req.Limit),
		},
	})
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
