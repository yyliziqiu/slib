package scsv

import (
	"testing"
)

type User struct {
	Name string `csv:"姓名"`
	Age  int
}

func TestSaveModels(t *testing.T) {
	users := []User{
		{"John", 30},
		{"Jane", 40},
	}

	err := SaveModels("/private/ws/self/slib/data/csv-test.csv", users)
	if err != nil {
		t.Error(err)
	}
}
