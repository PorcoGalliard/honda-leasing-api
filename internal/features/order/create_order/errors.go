package create_order

import "errors"

var (
	ErrMotorNotFound     = errors.New("motor tidak ditemukan")
	ErrMotorNotAvailable = errors.New("motor tidak tersedia (status bukan 'ready')")
	ErrDPExceedsPrice    = errors.New("DP tidak boleh melebihi atau sama dengan harga OTR motor")
	ErrNIKRequired       = errors.New("NIK wajib diisi untuk pendaftaran profil customer pertama kali")
)
