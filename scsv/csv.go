package scsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yyliziqiu/slib/sfile"
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

func SaveModels(filename string, models []any) error {
	if len(models) == 0 {
		return nil
	}
	rows := make([][]string, 0, len(models)+1)
	rows = append(rows, structFields(models[0]))
	for _, model := range models {
		rows = append(rows, structValues(model))
	}
	return Save(filename, rows)
}
