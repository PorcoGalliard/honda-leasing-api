package credit_simulation

import "errors"

var (
	ErrMotorNotFound  = errors.New("motor tidak ditemukan")
	ErrDPExceedsPrice = errors.New("DP tidak boleh melebihi atau sama dengan harga OTR motor")
)
