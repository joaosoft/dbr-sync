package json

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

var typeMarshal = reflect.TypeOf((*imarshal)(nil)).Elem()

type imarshal interface {
	MarshalJSON() ([]byte, error)
}

type marshal struct {
	object interface{}
	result *bytes.Buffer
	tags   []string
}

func newMarshal(object interface{}, tags ...string) *marshal {
	if len(tags) == 0 {
		tags = append(tags, jsonTag)
	}
	return &marshal{object: object, result: bytes.NewBuffer(make([]byte, 0)), tags: tags}
}

func (m *marshal) execute() ([]byte, error) {
	err := m.do(reflect.ValueOf(m.object))
	return m.result.Bytes(), err
}

func (m *marshal) do(object reflect.Value) error {
	types := reflect.TypeOf(object.Interface())

	handled, err := m.handleMarshalJSON(object)
	if err != nil {
		return err
	}

	if handled {
		return nil
	}

	if !object.CanInterface() {
		return nil
	}

	if object.Kind() == reflect.Ptr && !object.IsNil() {
		object = object.Elem()

		if object.IsValid() {
			types = object.Type()
		} else {
			return nil
		}
	}

	switch object.Kind() {
	case reflect.Struct:
		if _, err := m.result.WriteString(jsonStart); err != nil {
			return err
		}

		addComma := false
		for i := 0; i < types.NumField(); i++ {
			nextValue := object.Field(i)
			nextType := types.Field(i)

			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			handled, err := m.handleMarshalJSON(nextValue)
			if err != nil {
				return err
			}

			if handled {
				continue
			}

			if nextValue.Kind() == reflect.Ptr && !nextValue.IsNil() {
				nextValue = nextValue.Elem()
			}

			if !nextValue.CanInterface() {
				continue
			}

			exists, tag, err := m.loadTag(nextType)
			if err != nil {
				return err
			}

			if !exists {
				continue
			}

			if _, err := m.result.WriteString(fmt.Sprintf(`%s%s%s%s`, stringStartEnd, tag, stringStartEnd, is)); err != nil {
				return err
			}
			addComma = true

			if err := m.do(nextValue); err != nil {
				return err
			}
		}

		if _, err := m.result.WriteString(jsonEnd); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		if object.IsNil() {
			if _, err := m.result.WriteString(null); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(arrayStart); err != nil {
			return err
		}

		addComma := false
		for i := 0; i < object.Len(); i++ {

			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			nextValue := object.Index(i)

			handled, err := m.handleMarshalJSON(nextValue)
			if err != nil {
				return err
			}

			if handled {
				continue
			}

			if !nextValue.CanInterface() {
				continue
			}

			if err := m.do(nextValue); err != nil {
				return err
			}
			addComma = true
		}

		if _, err := m.result.WriteString(arrayEnd); err != nil {
			return err
		}

	case reflect.Map:
		if object.IsNil() {
			if _, err := m.result.WriteString(null); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(jsonStart); err != nil {
			return err
		}

		addComma := false
		for _, key := range object.MapKeys() {
			if addComma {
				if _, err := m.result.WriteString(comma); err != nil {
					return err
				}
			}

			nextValue := object.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if err := m.handleKey(key); err != nil {
				return err
			}

			if err := m.do(nextValue); err != nil {
				return err
			}
			addComma = true
		}

		if _, err := m.result.WriteString(jsonEnd); err != nil {
			return err
		}

	default:
		if err := m.handleValue(object); err != nil {
			return err
		}
	}
	return nil
}

func (m *marshal) handleKey(key reflect.Value) error {

	handled, err := m.handleMarshalJSON(key)
	if err != nil {
		return err
	}

	if handled {
		return nil
	}

	switch key.Kind() {
	case reflect.String:
		if _, err := m.result.WriteString(fmt.Sprintf("%s%s", m.encodeString(fmt.Sprintf(`%+v`, key.Interface())), is)); err != nil {
			return err
		}
	default:
		if _, err := m.result.WriteString(fmt.Sprintf(`%+v%s`, key.Interface(), is)); err != nil {
			return err
		}
	}

	return nil
}

func (m *marshal) handleValue(object reflect.Value) error {

	handled, err := m.handleMarshalJSON(object)
	if err != nil {
		return err
	}

	if handled {
		return nil
	}

	switch object.Kind() {
	case reflect.String:
		if _, err := m.result.WriteString(m.encodeString(fmt.Sprintf(`%+v`, object.Interface()))); err != nil {
			return err
		}
	default:
		if object.Kind() == reflect.Ptr && object.IsNil() {
			if _, err := m.result.WriteString(fmt.Sprintf(`%s`, null)); err != nil {
				return err
			}
			return nil
		}

		if _, err := m.result.WriteString(fmt.Sprintf(`%+v`, object.Interface())); err != nil {
			return err
		}
	}

	return nil
}

func (m *marshal) handleMarshalJSON(object reflect.Value) (bool, error) {
	val, ok := object.Interface().(imarshal)
	if ok {
		byts, err := val.MarshalJSON()
		if err != nil {
			return true, err
		}

		if _, err := m.result.WriteString(fmt.Sprintf(`%s`, string(byts))); err != nil {
			return true, err
		}
	}

	return ok, nil
}

func (m *marshal) loadTag(typ reflect.StructField) (exists bool, tag string, err error) {
	for _, searchTag := range m.tags {
		tag, exists = typ.Tag.Lookup(searchTag)

		if exists && tag == "-" {
			tag = ""
			exists = false
			break
		}

		if exists {
			break
		}
	}

	return exists, tag, err
}

func (m *marshal) encodeString(str string) string {
	return fmt.Sprintf("%s%s%s", stringStartEnd, strings.Replace(str, stringStartEnd, stringStartEndScaped, -1), stringStartEnd)
}
