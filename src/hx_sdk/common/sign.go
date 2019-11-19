package common

import (
	"encoding/hex"
	"fmt"
)

// der signature format
// https://crypto.stackexchange.com/questions/57731/ecdsa-signature-rs-to-asn1-der-encoding-question
func ConvertDerSig(sig string) ([]byte, error) {
	sigdata, err := hex.DecodeString(sig)
	if err != nil {
		return nil, err
	}
	if len(sigdata) != 70 {
		return nil, fmt.Errorf("der signature length should be 70, but %d", len(sigdata))
	}

	r := sigdata[4:36]
	s := sigdata[38:70]

	return append(r, s...), nil
}

func ConvertToDerSigB(bsig []byte) []byte {
	var r, s []byte

	r = append(r, 0)
	r = append(r, bsig[1:33]...)
	s = append(s, 0)
	s = append(s, bsig[33:65]...)
	for {
		if len(r) > 1 && r[0] == 0 && r[1] < 0x80 {
			r = r[1:]
		} else {
			break
		}
	}
	for {
		if len(s) > 1 && s[0] == 0 && s[1] < 0x80 {
			s = s[1:]
		} else {
			break
		}
	}

	size := 6 + len(r) + len(s)
	signedData := make([]byte, size, size)
	signedData[0] = 0x30
	signedData[1] = 4 + byte(len(r)) + byte(len(s))
	signedData[2] = 0x2
	signedData[3] = byte(len(r))
	copy(signedData[4:4+len(r)], r)
	signedData[4+len(r)] = 0x2
	signedData[5+len(r)] = byte(len(s))
	copy(signedData[6+len(r):6+len(r)+len(s)], s)

	return signedData

}

func ConvertToDerSig(sig string) ([]byte, error) {
	bsig, err := hex.DecodeString(sig)
	if err != nil {
		return nil, err
	}

	return ConvertToDerSigB(bsig), nil
}
