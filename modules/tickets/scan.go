package tickets

import (
	otp2 "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

func TokenToOTP(token string) (otp string, err error) {
	return totp.GenerateCodeCustom(token, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp2.DigitsSix,
		Algorithm: otp2.AlgorithmSHA1,
	})
}

func VerifyOTP(token string, otp string) (success bool, err error) {
	success, err = totp.ValidateCustom(token, otp, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp2.DigitsSix,
		Algorithm: otp2.AlgorithmSHA1,
	})

	return
}
