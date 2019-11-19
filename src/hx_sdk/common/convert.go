package common

import (
	"fmt"
	"math/big"
	"strings"
)

var _ = fmt.Printf

func AddZero(origin string, count int) string {
	for i := 0; i < count; i++ {
		origin += "0"
	}
	return origin
}

// 转换小数
func ConvertBigAmount(input string, percisions int) *big.Int {
	precision := big.NewInt(10)
	for i := 0; i < percisions-1; i++ {
		precision.Mul(precision, big.NewInt(10))
	}
	position := strings.Index(input, ".")
	var samt string
	if position == -1 {
		samt = AddZero(input, percisions)
	} else {
		res := input[:position]
		count := percisions - (len(input) - position - 1)
		res += input[position+1:]
		samt = AddZero(res, count)
		samt = strings.TrimLeft(samt, "0")
	}

	amt := &big.Int{}
	amt.SetString(samt, 10)
	return amt
}

func ConvertWithPrecision(s string, precision int64) (int64, error) {
	if precision > 20 {
		return 0, fmt.Errorf("param precision should be precision bits")
	}

	bf := big.Float{}
	_, _, err := bf.Parse(s, 10)
	if err != nil {
		return 0, err
	}

	fprec := convertPrecision(precision)
	f2 := bf.Mul(&bf, fprec)
	i2, _ := f2.Int64()

	return i2, nil
}

func convertPrecision(bit int64) *big.Float {
	prec := &big.Int{}
	x := big.NewInt(10)
	y := big.NewInt(bit)
	m := big.NewInt(0)

	prec.Exp(x, y, m)
	fprec := &big.Float{}
	fprec.SetInt(prec)

	// fmt.Println("precision:", fprec.String())
	return fprec
}

func ConvertToStringWithPrecision(s string, precision int64) (string, error) {
	/*

		bf := big.Float{}
		_, _, err := bf.Parse(s, 10)
		if err != nil {
			return "0", err
		}

		fmt.Println("param s:", bf.String())

		fprec := convertPrecision(precision)
		f2 := bf.Mul(&bf, fprec)

		bi, _ := f2.Int(nil)

		fmt.Println("to bigint:", bi.String())

		return bi.String(), nil
	*/

	if precision > 20 {
		return "0", fmt.Errorf("param precision should be precision bits")
	}
	bi := ConvertBigAmount(s, int(precision))
	return bi.String(), nil
}
