package common

import "crypto/sha256"

func HashSha256(data []byte) []byte {
	d := sha256.New()
	d.Write(data)
	return d.Sum(nil)
}
