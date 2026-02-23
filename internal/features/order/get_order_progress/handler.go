package get_order_progress

import (
	"net/http"
	"strconv"

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
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Tidak dapat mengidentifikasi user")
		return
	}

	contractIDStr := c.Param("contract_id")
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "contract_id tidak valid")
		return
	}

	contract, err := h.repo.FindContractByIDAndUserID(c.Request.Context(), contractID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data kontrak: "+err.Error())
		return
	}
	if contract == nil {
		response.Error(c, http.StatusNotFound, "Kontrak tidak ditemukan atau bukan milik Anda")
		return
	}

	tasks, err := h.repo.FindTasksByContractID(c.Request.Context(), contractID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil progress task: "+err.Error())
		return
	}

	taskItems := make([]ProgressTaskItem, 0, len(tasks))
	for _, t := range tasks {
		item := ProgressTaskItem{
			TaskID:      t.TaskID,
			TaskName:    t.TaskName,
			SequenceNo:  t.SequenceNo,
			Status:      t.Status,
			IsCompleted: t.Status == "completed",
		}
		if !t.Startdate.IsZero() {
			item.StartDate = t.Startdate.Format("02-Jan-2006")
		}
		if !t.ActualStartdate.IsZero() {
			item.ActualStartDate = t.ActualStartdate.Format("02-Jan-2006")
		}
		if !t.ActualEnddate.IsZero() {
			item.ActualEndDate = t.ActualEnddate.Format("02-Jan-2006")
		}
		taskItems = append(taskItems, item)
	}

	result := &OrderProgressResponse{
		ContractID:     contract.ContractID,
		ContractNumber: contract.ContractNumber,
		Status:         contract.Status,
		RequestDate:    contract.RequestDate.Format("02-Jan-2006"),
		Tasks:          taskItems,
	}

	response.Success(c, http.StatusOK, "Progress order berhasil diambil", result)
}
