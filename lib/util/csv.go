package util

import (
	"encoding/csv"
	"os"
)

func WriteCSV(path string, data [][]string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
