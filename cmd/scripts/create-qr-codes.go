package main

import (
	"fmt"

	"gitlab.com/markomalis/packback-api/pkg/services"
)

func main() {
	q := services.QRCodeService{}
	for i := 600; i <= 1400; i++ {
		err := q.WriteCode(
			fmt.Sprintf("https://packback.app/scanner?packback_id=%d", i),
			fmt.Sprintf("qr-codes/test-%d.png", i),
			256,
		)

		if err != nil {
			fmt.Println(err)
		}
	}
}
