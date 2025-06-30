package ssnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/sfile"
	"github.com/yyliziqiu/slib/slog"
)

type Snap struct {
	path   string
	data   any
	dupAge time.Duration
}

func New(path string, data interface{}) *Snap {
	return New2(path, data, 0)
}

func New2(path string, data interface{}, dupAge time.Duration) *Snap {
	return &Snap{path: path, data: data, dupAge: dupAge}
}

func New3(path string, data interface{}, poll time.Duration, n int) *Snap {
	return New2(path, data, poll*time.Duration(n)+time.Second)
}

func (s *Snap) Path() string {
	return s.path
}

func (s *Snap) Data() any {
	return s.data
}

func (s *Snap) DupAge() time.Duration {
	return s.dupAge
}

func (s *Snap) Save() error {
	return Save(s.path, s.data)
}

func (s *Snap) Load() error {
	return Load(s.path, s.data)
}

func (s *Snap) Duplicate() error {
	return Duplicate(s.path, s.data, s.dupAge)
}

func Save(path string, data interface{}) error {
	err := sfile.MakeDir(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("mkdir snap dir %s failed [%v]", filepath.Dir(path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snap data %s failed [%v]", path, err)
	}

	temp := path + ".temp"
	err = os.WriteFile(temp, bytes, 0644)
	if err != nil {
		return fmt.Errorf("save snap data %s failed [%v]", path, err)
	}
	err = os.Rename(temp, path)
	if err != nil {
		return fmt.Errorf("rename snap file %s failed [%v]", path, err)
	}

	return nil
}

func Load(path string, data interface{}) error {
	ok, err := sfile.Exist(path)
	if err != nil {
		return fmt.Errorf("check snap file %s failed [%v]", path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load snap file %s failed [%v]", path, err)
	}
	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, data)
}

func Duplicate(path string, data any, d time.Duration) error {
	if d <= 0 {
		return errors.New("duplication's age must be greater than zero")
	}

	nameRaw := filepath.Base(path)
	baseDir := filepath.Join(filepath.Dir(path), nameRaw+"-dup")

	// 清理过期快找副本
	_ = filepath.Walk(baseDir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Errorf("Walk snap duplicate failed, error:%v.", err)
			}
			return nil
		}
		if info.IsDir() || info.ModTime().After(time.Now().Add(-d)) {
			return nil
		}
		return os.Remove(file)
	})

	// 保存最新快照
	return Save(filepath.Join(baseDir, time.Now().Format("20060102-150405")), data)
}
