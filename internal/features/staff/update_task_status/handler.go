package update_task_status

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/internal/domain/query/account"
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

	contractID, err := strconv.ParseInt(c.Param("contract_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "contract_id tidak valid")
		return
	}

	taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "task_id tidak valid")
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	ctx := c.Request.Context()

	task, err := findTask(ctx, h.db, contractID, taskID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data task: "+err.Error())
		return
	}
	if task == nil {
		response.Error(c, http.StatusNotFound, ErrTaskNotFound.Error())
		return
	}

	if roleName != "SUPER_ADMIN" {
		roleID, err := getRoleIDByName(ctx, h.db, roleName)
		if err != nil || roleID != task.RoleID {
			response.Error(c, http.StatusForbidden, ErrTaskForbidden.Error())
			return
		}
	}

	oldStatus := task.Status
	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == "completed" {
		updates["actual_enddate"] = time.Now()
		if task.ActualStartdate.IsZero() {
			updates["actual_startdate"] = time.Now()
		}
	}

	lq := leasing.Use(h.db)
	lt := lq.LeasingTask
	_, err = lt.WithContext(ctx).
		Where(lt.TaskID.Eq(taskID)).
		Updates(updates)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal update task: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Status task berhasil diupdate", UpdateTaskStatusResponse{
		TaskID:     task.TaskID,
		TaskName:   task.TaskName,
		StatusLama: oldStatus,
		StatusBaru: req.Status,
	})
}

func findTask(ctx context.Context, db *gorm.DB, contractID, taskID int64) (*models.LeasingTask, error) {
	lq := leasing.Use(db)
	lt := lq.LeasingTask
	task, err := lt.WithContext(ctx).
		Where(lt.TaskID.Eq(taskID)).
		Where(lt.ContractID.Eq(contractID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func getRoleIDByName(ctx context.Context, db *gorm.DB, roleName string) (int64, error) {
	aq := account.Use(db)
	role, err := aq.Role.WithContext(ctx).
		Where(aq.Role.RoleName.Eq(roleName)).
		First()
	if err != nil {
		return 0, err
	}
	return role.RoleID, nil
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
