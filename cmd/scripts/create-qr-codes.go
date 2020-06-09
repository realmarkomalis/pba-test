package main

import (
	"fmt"

	"gitlab.com/markomalis/packback-api/pkg/services"
)

func main() {
	q := services.QRCodeService{}
<<<<<<< HEAD
	for i := 600; i <= 1400; i++ {
		err := q.WriteCode(
			fmt.Sprintf("https://packback.app/scanner?packback_id=%d", i),
=======
	for i := 1; i <= 500; i++ {
		_ = q.WriteCode(
			fmt.Sprintf("packback.app/?packback_id=%d", i),
>>>>>>> master
			fmt.Sprintf("qr-codes/test-%d.png", i),
			256,
		)

		if err != nil {
			fmt.Println(err)
		}
	}
}
