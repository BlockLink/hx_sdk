package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertEthAmount(t *testing.T) {
	var (
		ls = []string{
			"0.00123",
			"0.204308",
			"1.008",
			"3.1415926",
			"5",
			"15.319804",
			"66",
			"88",
			"309",
			"9999997.4444447901",
			"9999997.4432144031",
			"9999997.4432144132",
			"9999997.4432144233",
			"9999997.4432144344",
			"9999997.4432144445",
			"9999997.4432144546",
			"9999997.4432144647",
			"9999997.4432144748",
			"9999997.4432144859",
			"9999997.4432144955",

			"9999997.4432145011",
			"9999997.4432145132",
			"9999997.4432145293",
			"9999997.4432145344",
			"9999997.4432145465",
			"9999997.4432145536",
			"9999997.4432145647",
			"9999997.4432145758",
			"9999997.4432145817",
			"9999997.4432145929",
		}
	)

	for _, s := range ls {
		res := ConvertBigAmount(s, 18)
		// fmt.Println("ConvertBigAmount:", res.String())

		ret, err := ConvertWithPrecision(s, 18)
		assert.Nil(t, err)
		// fmt.Println("ConvertWithPrecision:", fmt.Sprint(ret))

		rets, err := ConvertToStringWithPrecision(s, 18)
		assert.Nil(t, err)
		// fmt.Println("ConvertToStringWithPrecision:", fmt.Sprint(rets))
		fmt.Println(s, ":", res, ret, rets)
	}
}
