package skvs

import (
	"time"
)

// ICGet deprecated
func (k Kvs) ICGet(key string) string {
	return k.Get(lower(key))
}

// ICS deprecated
func (k Kvs) ICS(key string, def string) string {
	return k.S(lower(key), def)
}

// ICB deprecated
func (k Kvs) ICB(key string, def bool) bool {
	return k.B(lower(key), def)
}

// ICI deprecated
func (k Kvs) ICI(key string, def int) int {
	return k.I(lower(key), def)
}

// ICI64 deprecated
func (k Kvs) ICI64(key string, def int64) int64 {
	return k.I64(lower(key), def)
}

// ICF64 deprecated
func (k Kvs) ICF64(key string, def float64) float64 {
	return k.F64(lower(key), def)
}

// ICD deprecated
func (k Kvs) ICD(key string, def time.Duration) time.Duration {
	return k.D(lower(key), def)
}
