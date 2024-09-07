package csvwriter

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"reflect"
)

type CSVWriter interface {
	CreateCSVFile(filename string, data interface{}) (string, error)
}

type CsvWriter struct {
	basicPath string
}

func NewCsvWriter(basicPath string) *CsvWriter {
	return &CsvWriter{basicPath: basicPath}
}

func (w *CsvWriter) CreateCSVFile(filename string, data interface{}) (string, error) {
	file, err := os.Create(path.Join(w.basicPath, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return "", fmt.Errorf("data must be a slice")
	}

	headers := w.getHeaders(value.Index(0).Interface())
	err = csvWriter.Write(headers)
	if err != nil {
		return "", err
	}

	for i := 0; i < value.Len(); i++ {
		record := w.getRecord(value.Index(i).Interface())
		err = csvWriter.Write(record)
		if err != nil {
			return "", err
		}
	}
	return filename, nil
}

func (w *CsvWriter) getHeaders(data interface{}) []string {
	var headers []string
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Struct {
		typ := value.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			headers = append(headers, field.Name)
		}
	}
	return headers
}

func (w *CsvWriter) getRecord(data interface{}) []string {
	var record []string
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Struct {
		typ := value.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := value.Field(i)
			record = append(record, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return record
}
