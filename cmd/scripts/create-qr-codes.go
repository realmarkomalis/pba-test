package main

import (
	"fmt"

	"gitlab.com/markomalis/packback-api/pkg/services"
)

func main() {
	q := services.QRCodeService{}
	for i := 1; i <= 600; i++ {
		_ = q.WriteCode(
			fmt.Sprintf("packback.app/?packback_id=%d", i),
			fmt.Sprintf("qr-codes/test-%d.png", i),
			256,
		)
	}
}
