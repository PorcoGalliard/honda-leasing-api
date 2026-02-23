package credit_simulation

type CreditSimulationRequest struct {
	MotorID int64   `json:"motor_id" binding:"required,min=1"`
	DP      float64 `json:"dp" binding:"required,min=0"`
}

type CreditSimulationResponse struct {
	MotorID   int64           `json:"motor_id"`
	NamaMotor string          `json:"nama_motor"`
	HargaOtr  float64         `json:"harga_otr"`
	DP        float64         `json:"dp"`
	Pokok     float64         `json:"pokok"`    // HargaOtr - DP
	Simulasi  []TenorSimulasi `json:"simulasi"` // Simulasi per tenor
}

type TenorSimulasi struct {
	Tenor            int     `json:"tenor"`
	TenorLabel       string  `json:"tenor_label"`        // ex.: "23 Bulan"
	AngsuranPerBulan float64 `json:"angsuran_per_bulan"` // saat ini saya mengimplementasikan bunga 0%
	TotalBayar       float64 `json:"total_bayar"`        // Pokok + DP (sama karena bunga 0%)
}
