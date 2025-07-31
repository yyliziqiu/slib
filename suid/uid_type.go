package uid

import (
	"errors"
	"strconv"
)

type Seed struct {
	A int64  // 当前时间戳
	B string // 当前时间戳16进制表示
	C int64  // 递增序号
}

var (
	ErrTimeBackForward = errors.New("time back forward")

	_padding = []string{"", "0", "00", "000", "0000", "00000", "000000", "0000000", "00000000", "000000000", "0000000000", "00000000000", "000000000000", "0000000000000", "00000000000000", "000000000000000", "0000000000000000"}

	_uid = New(1)
)

func hex(n int64, l int) string {
	s := strconv.FormatInt(n, 16)
	if len(s) < l {
		s = _padding[l-len(s)] + s
	}
	return s
}

func Get() string {
	return _uid.Get()
}

func GetOrFail() (string, error) {
	return _uid.GetOrFail()
}
