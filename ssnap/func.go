package ssnap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yyliziqiu/slib/sfile"
)

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
