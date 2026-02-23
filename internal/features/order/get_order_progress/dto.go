package get_order_progress

type OrderProgressResponse struct {
	ContractID     int64              `json:"contract_id"`
	ContractNumber string             `json:"contract_number"`
	Status         string             `json:"status"`
	RequestDate    string             `json:"request_date"`
	Tasks          []ProgressTaskItem `json:"tasks"`
}

type ProgressTaskItem struct {
	TaskID          int64  `json:"task_id"`
	TaskName        string `json:"task_name"`
	SequenceNo      int16  `json:"sequence_no"`
	Status          string `json:"status"`           
	IsCompleted     bool   `json:"is_completed"`     
	StartDate       string `json:"start_date"`       
	ActualStartDate string `json:"actual_start_date"`
	ActualEndDate   string `json:"actual_end_date"` 
}
