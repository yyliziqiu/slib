package ssnap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yyliziqiu/slib/sfile"
)

type Snap struct {
	Path string
	Data interface{}
}

func New(path string) *Snap {
	return &Snap{Path: path}
}

func New2(path string, data interface{}) *Snap {
	return &Snap{Path: path, Data: data}
}

func NewWithData(path string, data interface{}) *Snap {
	return New2(path, data)
}

func (s *Snap) Save() error {
	return s.SaveData(s.Data)
}

func (s *Snap) SaveData(data interface{}) error {
	err := sfile.MakeDir(filepath.Dir(s.Path))
	if err != nil {
		return fmt.Errorf("mkdir snap dir %s failed [%v]", filepath.Dir(s.Path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snap data %s failed [%v]", s.Path, err)
	}

	temp := ".temp"

	err = os.WriteFile(s.Path+temp, bytes, 0644)
	if err != nil {
		return fmt.Errorf("save snap data %s failed [%v]", s.Path, err)
	}

	err = os.Rename(s.Path+temp, s.Path)
	if err != nil {
		return fmt.Errorf("rename snap file %s failed [%v]", s.Path, err)
	}

	return nil
}

func (s *Snap) Load() error {
	return s.LoadData(s.Data)
}

func (s *Snap) LoadData(data interface{}) error {
	ok, err := sfile.Exist(s.Path)
	if err != nil {
		return fmt.Errorf("check snap file %s failed [%v]", s.Path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("load snap file %s failed [%v]", s.Path, err)
	}

	return json.Unmarshal(bytes, data)
}
