package impt

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetIndexesByTag(modelType reflect.Type, tagName string) (map[int]string, error) {
	ma := make(map[int]string, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		if len(tagValue) > 0 {
			ma[i] = tagValue
		} else {
			ma[i] = ""
		}
	}
	return ma, nil
}
func NewCSVFormatter(modelType reflect.Type) *CSVFormatter {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		panic("error get formatCols")
	}
	return &CSVFormatter{modelType: modelType, formatCols: formatCols}
}

type CSVFormatter struct {
	modelType  reflect.Type
	formatCols map[int]string
}

func (f CSVFormatter) ToStruct(ctx context.Context, lines []string) (interface{}, error) {
	record := reflect.New(f.modelType).Interface()
	err := ScanLine(lines, f.modelType, &record, f.formatCols)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return reflect.Indirect(reflect.ValueOf(record)).Interface(), nil
	}
	return record, err
}

func ScanLine(lines []string, modelType reflect.Type, record interface{}, formatCols map[int]string) error {
	s := reflect.Indirect(reflect.ValueOf(record))
	numFields := s.NumField()
	for i := 0; i < numFields; i++ {
		field := modelType.Field(i)
		typef := field.Type.String()
		line := lines[i]
		f := s.Field(i)
		if f.CanSet() {
			switch typef {
			case "string", "*string":
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&line))
				} else {
					f.SetString(line)
				}
			case "int64", "*int64":
				value, _ := strconv.ParseInt(line, 10, 64)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&value))
				} else {
					f.SetInt(value)
				}
			case "int", "*int":
				value, _ := strconv.Atoi(line)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&value))
				} else {
					f.Set(reflect.ValueOf(value))
				}
			case "bool":
				boolValue, _ := strconv.ParseBool(line)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&boolValue))
				} else {
					f.SetBool(boolValue)
				}
			case "float64", "*float64":
				floatValue, _ := strconv.ParseFloat(line, 64)
				if f.Kind() == reflect.Ptr {
					f.Set(reflect.ValueOf(&floatValue))
				} else {
					f.SetFloat(floatValue)
				}
			case "time.Time", "*time.Time":
				if format, ok := formatCols[i]; ok {
					if strings.Contains(format, "dateFormat:") {
						layoutDateStr := strings.ReplaceAll(format, "dateFormat:", "")
						fieldDate, err := time.Parse(layoutDateStr, line)
						if err != nil {
							return err
						}
						if f.Kind() == reflect.Ptr {
							f.Set(reflect.ValueOf(&fieldDate))
						} else {
							f.Set(reflect.ValueOf(fieldDate))
						}
					}
				}
			}
		}
	}
	return nil
}