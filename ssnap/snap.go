package ssnap

import (
	"path/filepath"
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

func (s *Snap) SaveDuplicate(n int) error {
	name := filepath.Base(s.path)
	path := filepath.Join(filepath.Dir(s.path), "duplicates", name)
	return Save(path, s.data)
}
