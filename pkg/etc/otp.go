package etc

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateOTP(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(otp)
}
