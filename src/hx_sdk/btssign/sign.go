package btssign

import (
	secp "github.com/bitnexty/secp256k1-go"
)

func SignCompact(msg []byte, seckey []byte, requireCanonical bool) ([]byte, error) {
	// BtsSign has check IsCanonical
	return secp.BtsSign(msg, seckey, requireCanonical)
}

func IsCanonical(sig []byte) bool {
	tmp := (sig[0]&0x80 != 0) || (sig[0] == 0x0 && (sig[1]&0x80 != 0)) || (sig[32]&0x80 != 0) || (sig[32] == 0x0 && (sig[33]&0x80 != 0))

	return !tmp
}

func IsCanonicalv2(sig []byte) bool {
	return !(sig[0]&0x80 != 0) && !(sig[0] == 0 && (sig[1]&0x80) == 0) && (sig[32]&0x80 == 0) && !(sig[32] == 0x0 && (sig[33]&0x80 != 0))
}
