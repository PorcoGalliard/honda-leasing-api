package update_order_status

import "errors"

var validTransitions = map[string]map[string]bool{
	"draft":    {"approved": true, "canceled": true},
	"approved": {"active": true, "canceled": true},
	"active":   {"late": true, "paid_off": true, "repo": true, "canceled": true},
	"late":     {"active": true, "paid_off": true, "repo": true, "canceled": true},
}

var roleTransitionAllowed = map[string]map[string]bool{
	"SUPER_ADMIN":  {"approved": true, "active": true, "late": true, "paid_off": true, "repo": true, "canceled": true},
	"ADMIN_CABANG": {"approved": true, "active": true, "late": true, "paid_off": true, "repo": true, "canceled": true},
	"FINANCE":      {"active": true, "paid_off": true},
	"COLLECTION":   {"late": true},
}

var (
	ErrContractNotFound    = errors.New("kontrak tidak ditemukan")
	ErrInvalidTransition   = errors.New("perpindahan status tidak valid")
	ErrForbiddenTransition = errors.New("role Anda tidak diizinkan melakukan perubahan status ini")
)

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved active late paid_off repo canceled"`
}

type UpdateOrderStatusResponse struct {
	ContractID     int64  `json:"contract_id"`
	ContractNumber string `json:"contract_number"`
	StatusLama     string `json:"status_lama"`
	StatusBaru     string `json:"status_baru"`
}
