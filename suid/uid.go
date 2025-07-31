package uid

import (
	"sync"
	"time"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/ssnap"
)

type Uid struct {
	node string
	seed Seed
	snap *ssnap.Snap
	mu   sync.Mutex
}

func New(node int) *Uid {
	return New2(node, "")
}

func New2(node int, path string) *Uid {
	t := &Uid{
		node: hex(int64(node), 2),
	}
	if path == "" {
		t.snap = ssnap.New(path, &t.seed)
	}

	return t
}

func (t *Uid) Save(_ bool) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.snap.Save()
}

func (t *Uid) Load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.snap.Load()
}

// Get 返回十六位唯一 ID
func (t *Uid) Get() string {
	id, err := t.GetOrFail()
	if err != nil {
		slog.Error("[Uuid.Get] Time back forward.")
	}
	return id
}

// GetOrFail 返回十六位唯一 ID
func (t *Uid) GetOrFail() (string, error) {
	t.mu.Lock()

	nano := time.Now().UnixNano()
	curr := nano / 1e9
	if curr < t.seed.A {
		t.mu.Unlock()
		return "", ErrTimeBackForward
	}

	if curr > t.seed.A {
		t.seed.A = curr
		t.seed.B = hex(curr, 8)
		t.seed.C = (nano % 1e7) + 1048576 // 1048576 = 0x100000, 确保C转化为16进制后为6位数，且有一定的增长空间
	}
	t.seed.C++

	id := t.seed.B + t.node + hex(t.seed.C, 6)

	t.mu.Unlock()

	return id, nil
}
