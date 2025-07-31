package uid

var _padding = []string{
	"", "0", "00", "000", "0000", "00000", "000000", "0000000", "00000000", "000000000", "0000000000",
	"00000000000", "000000000000", "0000000000000", "00000000000000", "000000000000000", "0000000000000000",
}

type Seed struct {
	A int64  // 当前时间戳
	B string // 当前时间戳16进制表示
	C int64  // 递增序号
}

var _uid = New(1)

func Get() string {
	return _uid.Get()
}

func GetSeed() Seed {
	return _uid.seed
}

func SetSeed(seed Seed) {
	_uid.seed = seed
}
