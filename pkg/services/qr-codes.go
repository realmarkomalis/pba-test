package services

import (
	b64 "encoding/base64"

	qrcode "github.com/skip2/go-qrcode"
)

type IQRCodeSerive interface {
	CreateCode(size int) ([]byte, error)
	WriteCode(data, filename string, size int) error
}

type QRCodeService struct{}

func (q QRCodeService) CreateCode(data string, size int) ([]byte, error) {
	code, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func (q QRCodeService) CreateBase64Code(data string) (string, error) {
	code, err := q.CreateCode(data, 256)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(code), nil
}

func (q QRCodeService) WriteCode(data, filename string, size int) error {
	err := qrcode.WriteFile(data, qrcode.Medium, size, filename)
	if err != nil {
		return err
	}

	return nil
}
