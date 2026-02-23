package list_orders

type ListOrdersRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status string `form:"status" binding:"omitempty,oneof=draft approved active late paid_off repo canceled"`
}

type ListOrdersResponse struct {
	Orders     []OrderItem `json:"orders"`
	Pagination Pagination  `json:"pagination"`
}

type OrderItem struct {
	ContractID      int64   `json:"contract_id"`
	ContractNumber  string  `json:"contract_number"`
	RequestDate     string  `json:"request_date"`
	Status          string  `json:"status"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
	MotorMerk       string  `json:"motor_merk"`
	MotorType       string  `json:"motor_type"`
	NilaiKendaraan  float64 `json:"nilai_kendaraan"`
	DpDibayar       float64 `json:"dp_dibayar"`
	TenorBulan      int16   `json:"tenor_bulan"`
	CicilanPerBulan float64 `json:"cicilan_per_bulan"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}
