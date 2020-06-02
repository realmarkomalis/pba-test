package services

import qrcode "github.com/skip2/go-qrcode"

type IQRCodeSerive interface {
	CreateCode(size int) ([]byte, error)
	WriteCode(data, filename string, size int) error
}

type QRCodeService struct{}

func (q QRCodeService) CreateCode(size int) ([]byte, error) {
	return []byte{}, nil
}

func (q QRCodeService) WriteCode(data, filename string, size int) error {
	err := qrcode.WriteFile(data, qrcode.Medium, size, filename)
	if err != nil {
		return err
	}

	return nil
}
