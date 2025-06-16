package ssnap

import (
	"os"
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/slog"
)

type Snap struct {
	path string
	data any
}

func New(path string, data interface{}) *Snap {
	return &Snap{path: path, data: data}
}

func (s *Snap) Path() string {
	return s.path
}

func (s *Snap) Data() any {
	return s.data
}

func (s *Snap) Save() error {
	return Save(s.path, s.data)
}

func (s *Snap) Load() error {
	return Load(s.path, s.data)
}

func (s *Snap) SaveDuplicate(d time.Duration) error {
	name := filepath.Base(s.path)
	path := filepath.Join(filepath.Dir(s.path), name+"-dup")

	// 清理过期快找副本
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Errorf("Walk snap duplicate failed, error:%v.", err)
			}
			return nil

		}
		if info.IsDir() || info.ModTime().After(time.Now().Add(-d)) {
			return nil
		}
		return os.Remove(path)
	})

	// 保存最新快照
	return Save(filepath.Join(path, time.Now().Format(time.DateTime)), s.data)
}
