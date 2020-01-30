package dbr

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/joaosoft/errors"
)

var typeScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func GetEnv() string {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	return env
}

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(file string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(file) {
		return nil, errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(file string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(file) {
		return nil, errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(file string, obj interface{}) error {
	if !Exists(file) {
		return errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(file, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

func read(columns []interface{}, rows *sql.Rows, value reflect.Value) (int, error) {

	value = value.Elem() // get the addressable value
	isScanner := value.Addr().Type().Implements(typeScanner)
	isSlice := value.Kind() == reflect.Slice && !isScanner
	isMap := value.Kind() == reflect.Map && !isScanner
	isMapOfSlices := isMap && value.Type().Elem().Kind() == reflect.Slice

	count := 0

	// load each row
	for rows.Next() {
		var elem reflect.Value
		if isMapOfSlices {
			elem = reflectAlloc(value.Type().Elem().Elem())
		} else if isSlice || isMap {
			elem = reflectAlloc(value.Type().Elem())
		} else {
			elem = value
		}

		// load expression values
		fields, err := getFields(constLoadOptionRead, columns, elem)
		if err != nil {
			return 0, err
		}

		// scan values from row
		err = rows.Scan(fields...)
		if err != nil {
			return 0, err
		}

		count++
		if isSlice {
			value.Set(reflect.Append(value, elem))
		} else {
			break
		}
	}

	return count, nil
}

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}

func getFields(loadOption loadOption, columns []interface{}, object reflect.Value) ([]interface{}, error) {
	var fields []interface{}

	// add columns to a map
	mapColumns := make(map[interface{}]bool)
	for _, name := range columns {
		mapColumns[fmt.Sprint(name)] = true
	}

	mappedValues := make(map[interface{}]interface{})
	loadColumnStructValues(loadOption, columns, mapColumns, object, mappedValues)

	for _, name := range columns {
		if value, ok := mappedValues[fmt.Sprint(name)]; ok {
			fields = append(fields, value)
		} else {
			var value interface{}
			fields = append(fields, &value)
		}
	}

	return fields, nil
}

func loadColumnStructValues(loadOption loadOption, columns []interface{}, mapColumns map[interface{}]bool, object reflect.Value, mappedValues map[interface{}]interface{}) {

	if object.Kind() == reflect.Ptr && object.IsNil() {
		object.Set(reflect.New(object.Type().Elem()))
	}

	if object.CanAddr() && object.Addr().Type().Implements(typeScanner) {
		mappedValues[fmt.Sprint(columns[0])] = object.Addr().Interface()
		return
	}

	switch object.Kind() {
	case reflect.Ptr:
		loadColumnStructValues(loadOption, columns, mapColumns, object.Elem(), mappedValues)

	case reflect.Struct:
		t := object.Type()
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			field := object.Field(i)

			if structField.PkgPath != "" && !structField.Anonymous {
				// unexported
				continue
			}
			tag := structField.Tag.Get(string(loadOption))
			if tag == "-" {
				// ignore
				continue
			}

			if tag == "" {
				tag = structField.Tag.Get(string(constLoadOptionDefault))
				if tag == "-" {
					// ignore
					continue
				}
			}

			if tag == "" {
				switch field.Kind() {
				case reflect.Struct:
					loadColumnStructValues(loadOption, columns, mapColumns, field, mappedValues)
				default:
					continue
				}
			}

			isSlice := field.Kind() == reflect.Slice
			isMap := field.Kind() == reflect.Map
			isMapOfSlices := isMap && field.Type().Elem().Kind() == reflect.Slice

			if isMapOfSlices {
				field = reflectAlloc(field.Type().Elem())
			} else if isSlice || isMap {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			}

			if field.Kind() == reflect.Ptr && field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}

			if field.CanAddr() && field.Addr().Type().Implements(typeScanner) {
				mappedValues[tag] = field.Addr().Interface()
				continue
			}

			if _, ok := mapColumns[tag]; ok {
				mappedValues[tag] = field.Addr().Interface()
				continue
			}
		}

	default:
		mappedValues[fmt.Sprint(columns[0])] = object.Addr().Interface()
	}
}

func loadStructValues(loadOption loadOption, object reflect.Value, columns *[]interface{}, mappedValues map[interface{}]reflect.Value) {
	switch object.Kind() {
	case reflect.Ptr:
		if !object.IsNil() {
			loadStructValues(loadOption, object.Elem(), columns, mappedValues)
		}
	case reflect.Struct:
		t := object.Type()

		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			field := object.Field(i)

			if structField.PkgPath != "" && !structField.Anonymous {
				// unexported
				continue
			}

			tag := structField.Tag.Get(string(loadOption))
			if tag == "-" {
				// ignore
				continue
			}

			if tag == "" {
				tag = structField.Tag.Get(string(constLoadOptionDefault))
				if tag == "-" || tag == "" {
					// ignore
					continue
				}
			}

			if _, ok := mappedValues[tag]; !ok {
				mappedValues[tag] = field
				if columns != nil {
					*columns = append(*columns, tag)
				}
			}
		}
	}
}
