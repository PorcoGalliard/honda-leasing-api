package list_motors

type ListMotorsRequest struct {
	Page        int     `form:"page" binding:"omitempty,min=1"`
	Limit       int     `form:"limit" binding:"omitempty,min=1,max=100"`
	MotorType   string  `form:"motor_type"` // Filter berdasarkan motor type
	Status      string  `form:"status"`     // Filter berdasarkan status
	MinPrice    float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice    float64 `form:"max_price" binding:"omitempty,min=0"`
	GroupByType bool    `form:"group_by_type"`                                                     // Group results by motor type
	SortBy      string  `form:"sort_by" binding:"omitempty,oneof=motor_type harga_otr created_at"` // Sort field
	OrderBy     string  `form:"order_by" binding:"omitempty,oneof=asc desc"`                       // asc atau desc
}

// Response untuk list motor dengan grouping
type ListMotorsResponse struct {
	MotorsByType []MotorTypeGroup `json:"motors_by_type,omitempty"` // Jika group_by_type=true
	Motors       []MotorItem      `json:"motors,omitempty"`         // Jika group_by_type=false
	Pagination   *Pagination      `json:"pagination,omitempty"`     // Hanya ada jika tidak grouping
	TotalMotors  int              `json:"total_motors"`             // Total semua motor
}

// Grouping motor berdasarkan type
type MotorTypeGroup struct {
	MotorType     string      `json:"motor_type"`      // Classic, Sport, Matic, Maxi, Bebek
	MotorCount    int         `json:"motor_count"`     // Jumlah motor dalam type ini
	Motors        []MotorItem `json:"motors"`
}

type MotorItem struct {
	MotorID       int64    `json:"motor_id"`
	Merk          string   `json:"merk"`
	MotorType     string   `json:"motor_type"`
	Tahun         int16    `json:"tahun"`
	Warna         string   `json:"warna"`
	NomorRangka   string   `json:"nomor_rangka"`
	NomorMesin    string   `json:"nomor_mesin"`
	CcMesin       string   `json:"cc_mesin"`
	NomorPolisi   string   `json:"nomor_polisi"`
	StatusUnit    string   `json:"status_unit"`
	HargaOtr      float64  `json:"harga_otr"`
	Images        []string `json:"images"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}
