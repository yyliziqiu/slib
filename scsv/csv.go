package scsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/yyliziqiu/slib/sfile"
	"github.com/yyliziqiu/slib/sreflect"
)

func Save(filename string, rows [][]string) error {
	// 创建存储目录
	err := sfile.MakeDir(filepath.Dir(filename))
	if err != nil {
		return fmt.Errorf("mkdir failed [%v]", err)
	}

	// 优化文件名
	if !strings.HasSuffix(filename, ".csv") {
		filename = filename + ".csv"
	}

	// 创建 CSV 文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create CSV file failed [%v]", err)
	}
	defer file.Close()

	// 写入 CSV 文件
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(rows)
	if err != nil {
		return fmt.Errorf("write date to CSV failed [%v]", err)
	}

	return nil
}

func SaveModels(filename string, models any) error {
	v := reflect.ValueOf(models)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("modeles type must be slice or array")
	}

	size := v.Len()
	if size == 0 {
		return nil
	}

	head := v.Index(0).Interface()

	rows := make([][]string, 0, size+1)
	rows = append(rows, titles(head))
	for i := 0; i < size; i++ {
		rows = append(rows, sreflect.ValuesOf(v.Index(i).Interface()))
	}

	return Save(filename, rows)
}

func titles(s any) []string {
	mt := reflect.TypeOf(s)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		title := mt.Field(i).Tag.Get("csv")
		if title == "" {
			title = mt.Field(i).Name
		}
		fields = append(fields, title)
	}
	return fields
}
