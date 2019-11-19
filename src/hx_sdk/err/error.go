package err

// 封装 error, 让客户端更容易识别

import (
	"encoding/json"
	"fmt"
)

type ErrRet struct {
	ErrCode string
	ErrMsg  string
}

func ErrCode(code string, msg ...string) error {
	if msg == nil {
		return fmt.Errorf(`{"ErrCode": "%s"}`, code)
	}
	er := ErrRet{
		ErrCode: code,
		ErrMsg:  msg[0],
	}
	buf, _ := json.Marshal(&er)
	return fmt.Errorf(string(buf))
}

// if err is ErrRet format, return err;
// else wrap to ErrRet format
func ErrWrap(err error, code ...string) error {
	var er ErrRet

	if err == nil {
		return nil
	}

	e := json.Unmarshal([]byte(err.Error()), &er)
	if e == nil && er.ErrCode != "" {
		return err
	}

	if e != nil {
		if code != nil {
			er.ErrCode = code[0]
		} else {
			er.ErrCode = ErrCodeInternalError
		}
		if len(code) > 1 {
			er.ErrMsg = code[1]
		} else {
			er.ErrMsg = err.Error()
		}
		buf, _ := json.Marshal(&er)
		return fmt.Errorf(string(buf))
	}

	// is ErrRet format
	// er.ErrCode == ""
	if code != nil {
		er.ErrCode = code[0]
	} else {
		er.ErrCode = ErrCodeInternalError
	}
	if er.ErrMsg == "" {
		if len(code) > 1 {
			er.ErrMsg = code[1]
		} else {
			er.ErrMsg = "Internal Error: " + er.ErrCode
		}
	}

	buf, _ := json.Marshal(&er)
	return fmt.Errorf(string(buf))
}
