package skvs

import (
	"strings"
	"time"

	"github.com/yyliziqiu/slib/sconv"
)

type Kvs map[string]string

// 1

func (k Kvs) String(key string) (string, bool) {
	if val, ok := k[key]; ok {
		return strings.TrimSpace(val), true
	}
	return "", false
}

func (k Kvs) Bool(key string) (bool, bool) {
	if val, ok := k.String(key); ok {
		return sconv.S2B(val), true
	}
	return false, false
}

func (k Kvs) Int(key string) (int, bool) {
	if val, ok := k.String(key); ok {
		return sconv.S2I(val), true
	}
	return 0, false
}

func (k Kvs) Int64(key string) (int64, bool) {
	if val, ok := k.String(key); ok {
		return sconv.S2I64(val), true
	}
	return 0, false
}

func (k Kvs) Float64(key string) (float64, bool) {
	if val, ok := k.String(key); ok {
		return sconv.S2F64(val), true
	}
	return 0, false
}

func (k Kvs) Duration(key string) (time.Duration, bool) {
	if val, ok := k.String(key); ok {
		d, _ := time.ParseDuration(val)
		return d, true
	}
	return 0, false
}

// 2

func (k Kvs) S(key string, def string) string {
	if val, ok := k.String(key); ok {
		return val
	}
	return def
}

func (k Kvs) B(key string, def bool) bool {
	if val, ok := k.Bool(key); ok {
		return val
	}
	return def
}

func (k Kvs) I(key string, def int) int {
	if val, ok := k.Int(key); ok {
		return val
	}
	return def
}

func (k Kvs) I64(key string, def int64) int64 {
	if val, ok := k.Int64(key); ok {
		return val
	}
	return def
}

func (k Kvs) F64(key string, def float64) float64 {
	if val, ok := k.Float64(key); ok {
		return val
	}
	return def
}

func (k Kvs) D(key string, def time.Duration) time.Duration {
	if val, ok := k.Duration(key); ok {
		return val
	}
	return def
}

// 3

var lower = strings.ToLower

func (k Kvs) LS(key string, def string) string {
	return k.S(lower(key), def)
}

func (k Kvs) LB(key string, def bool) bool {
	return k.B(lower(key), def)
}

func (k Kvs) LI(key string, def int) int {
	return k.I(lower(key), def)
}

func (k Kvs) LI64(key string, def int64) int64 {
	return k.I64(lower(key), def)
}

func (k Kvs) LF64(key string, def float64) float64 {
	return k.F64(lower(key), def)
}

func (k Kvs) LD(key string, def time.Duration) time.Duration {
	return k.D(lower(key), def)
}

// 4

func (k Kvs) Get(key string) string {
	return k.S(key, "")
}

func (k Kvs) LGet(key string) string {
	return k.Get(lower(key))
}

func (k Kvs) Id() string {
	return k.S("id", "")
}

func (k Kvs) Name() string {
	return k.S("name", "")
}

func (k Kvs) Slug() string {
	return k.S("slug", "")
}

func (k Kvs) Type() string {
	return k.S("type", "")
}

func (k Kvs) Url() string {
	return k.S("url", "")
}

func (k Kvs) Key() string {
	return k.S("key", "")
}

func (k Kvs) Token() string {
	return k.S("token", "")
}

func (k Kvs) Enabled() bool {
	return k.B("enabled", false)
}

func (k Kvs) Disabled() bool {
	return k.B("disabled", false)
}
