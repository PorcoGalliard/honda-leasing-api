package create_order

type CreateOrderRequest struct {
	MotorID       int64   `json:"motor_id" binding:"required,min=1"`
	ContactName   string  `json:"contact_name" binding:"required,min=3,max=100"`
	PhoneNumber   string  `json:"phone_number" binding:"required"`
	NIK           string  `json:"nik" binding:"omitempty,len=16"` 
	DP            float64 `json:"dp" binding:"required,min=0"`
	Tenor         int16   `json:"tenor" binding:"required,oneof=23 29 35"`
	RequestDate   string  `json:"request_date" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	PromoDiscount float64 `json:"promo_discount" binding:"omitempty,min=0"`
	PromoName     string  `json:"promo_name" binding:"omitempty"`
}

type CreateOrderResponse struct {
	ContractID     int64           `json:"contract_id"`
	ContractNumber string          `json:"contract_number"`
	Status         string          `json:"status"`
	Motor          OrderMotorInfo  `json:"motor"`
	Tenor          int16           `json:"tenor"`
	TotalSummary   OrderSummary    `json:"total_summary"`
	Tasks          []OrderTaskItem `json:"tasks"`
}

// OrderMotorInfo - Info motor di dalam response order
type OrderMotorInfo struct {
	MotorID   int64   `json:"motor_id"`
	Merk      string  `json:"merk"`
	MotorType string  `json:"motor_type"`
	Tahun     int16   `json:"tahun"`
	HargaOtr  float64 `json:"harga_otr"`
}

// OrderSummary - Rincian kalkulasi biaya
type OrderSummary struct {
	HargaOtr         float64 `json:"harga_otr"`
	DP               float64 `json:"dp"`
	PromoDiscount    float64 `json:"promo_discount"`
	BiayaAdmin       float64 `json:"biaya_admin"`
	Asuransi         float64 `json:"asuransi"`
	Fidusia          float64 `json:"fidusia"`
	Materai          float64 `json:"materai"`
	PokokPinjaman    float64 `json:"pokok_pinjaman"`
	CicilanPerBulan  float64 `json:"cicilan_per_bulan"`
	SubTotal         float64 `json:"subtotal"` // total yang harus dibayar customer (pokok + fees)
}

// OrderTaskItem - Item task di progress screen
type OrderTaskItem struct {
	TaskID     int64  `json:"task_id"`
	TaskName   string `json:"task_name"`
	SequenceNo int16  `json:"sequence_no"`
	Status     string `json:"status"` // inprogress, completed, cancelled
}
