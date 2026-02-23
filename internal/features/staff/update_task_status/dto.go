package update_task_status

import "errors"

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=completed cancelled"`
}

type UpdateTaskStatusResponse struct {
	TaskID     int64  `json:"task_id"`
	TaskName   string `json:"task_name"`
	StatusLama string `json:"status_lama"`
	StatusBaru string `json:"status_baru"`
}

var (
	ErrTaskNotFound  = errors.New("task tidak ditemukan")
	ErrTaskForbidden = errors.New("role Anda tidak memiliki izin untuk mengupdate task ini")
)
