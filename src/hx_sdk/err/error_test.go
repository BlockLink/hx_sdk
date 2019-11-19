package err

import (
	"fmt"
	"testing"
)

func TestErrCode(t *testing.T) {
	fmt.Println(ErrCode(ErrCodeInvalidAPI, "error: \"unsupport api schem\""))
}
