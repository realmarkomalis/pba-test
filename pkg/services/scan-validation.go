package services

import (
	"crypto/sha256"
	"fmt"
)

type IScanValidationService interface {
	Generate(id string) string
	Validate(id, token string) bool
}

type ScanValidationService struct {
	secretKey string
}

func NewScanValidationService(key string) *ScanValidationService {
	return &ScanValidationService{
		secretKey: key,
	}
}

func (s ScanValidationService) Generate(id string) string {
	h := sha256.New()
	h.Write(formatedBytes(id, s.secretKey))
	return string(h.Sum(nil))
}

func (s ScanValidationService) Validate(id, token string) bool {
	return token == s.Generate(id)
}

func formatedBytes(id, key string) []byte {
	return []byte(fmt.Sprintf(
		"%s/%s",
		id,
		key,
	))
}
