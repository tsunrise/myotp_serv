package test

import (
	"fmt"
	"myotp_serv/modules/tickets"
	"testing"
	"time"
)

func TestOTP(t *testing.T) {
	token := "WBEIHZLKCVVI3ZLDV6SGPVE55RVEB2NB2AD7AM6CQIYGDPOQ5KSZ7DSNUZBQCBJCBUFSS2ELONFY5IHT"

	for {
		otp, _ := tickets.TokenToOTP(token)

		fmt.Printf("%v -> %v\n", token, otp)

		time.Sleep(5 * time.Second)
	}
}
