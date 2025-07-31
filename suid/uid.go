package uid

import (
	"strconv"
	"sync"
	"time"

	"github.com/yyliziqiu/slib/slog"
)

type Uid struct {
	node string
	seed Seed
	mu   sync.Mutex
}

func New(node int) *Uid {
	t := &Uid{}
	t.node = t.hex(int64(node), 2)
	return t
}

// Get 返回十六位唯一 ID
func (t *Uid) Get() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	nano := time.Now().UnixNano()
	curr := nano / 1e9
	if curr < t.seed.A {
		slog.Error("[Uuid.Get] Time back forward.")
	}
	if curr > t.seed.A {
		t.seed.A = curr
		t.seed.B = t.hex(curr, 8)
		t.seed.C = (nano % 1e7) + 1048576 // 1048576 = 0x100000, 确保C转化为16进制后为6位数，且有一定的增长空间
	}

	t.seed.C++

	id := t.seed.B + t.node + t.hex(t.seed.C, 6)

	return id
}

func (t *Uid) hex(n int64, l int) string {
	s := strconv.FormatInt(n, 16)
	if len(s) < l {
		s = _padding[l-len(s)] + s
	}
	return s
}
