package validator

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func (v *Validator) _convertToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%+v", value)
}

func (v *Validator) _getValue(value reflect.Value) (isNil bool, _ reflect.Value, _ interface{}) {
again:
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return true, value, value.Interface()
		}
		value = value.Elem()

		goto again
	}

	return value.Interface() == nil, value, value.Interface()
}

func (v *Validator) _loadExpectedValue(context *ValidatorContext, expected interface{}) (interface{}, error) {

	if expected != nil && v._convertToString(expected) != "" {
		strValue := v._convertToString(expected)
		if matched, err := regexp.MatchString(constRegexForReplaceId, strValue); err != nil {
			return "", err
		} else {
			if matched {
				replacer := strings.NewReplacer(constTagReplaceIdStart, "", constTagReplaceIdEnd, "")
				id := replacer.Replace(strValue)

				value, ok := context.GetValue(constTagId, id)
				if !ok {
					value, ok = context.GetValue(constTagArg, id)
					if !ok {
						value, ok = context.GetValue(constTagJson, id)
					}
				}

				if ok {
					return value.value.Interface(), nil
				}
			}
		}

	}
	return expected, nil
}

func (v *Validator) _random(strValue string) string {
	rand.Seed(time.Now().UnixNano())
	alphabetLowerChars := []rune(constAlphanumericLowerAlphabet)
	alphabetUpperChars := []rune(constAlphanumericUpperAlphabet)
	alphabetNumbers := []rune(constNumericAlphabet)
	alphabetSpecial := []rune(constSpecialAlphabet)

	newValue := []rune(strValue)

	for i, char := range newValue {
		if !unicode.IsSpace(char) {
			var alphabet []rune
			if unicode.IsLetter(char) {
				if unicode.IsUpper(char) {
					alphabet = alphabetUpperChars
				} else {
					alphabet = alphabetLowerChars
				}
			} else if unicode.IsNumber(char) {
				alphabet = alphabetNumbers
			} else {
				alphabet = alphabetSpecial
			}

			newValue[i] = alphabet[rand.Intn(len(alphabet))]
		}
	}

	return string(newValue)
}

func _setValue(kind reflect.Kind, obj reflect.Value, newValue interface{}) (err error) {
	switch value := newValue.(type) {
	case string:
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var v int
			if v, err = strconv.Atoi(value); err != nil {
				return err
			}
			obj.SetInt(int64(v))
		case reflect.Float32, reflect.Float64:
			var v float64
			if v, err = strconv.ParseFloat(value, 64); err != nil {
				return err
			}
			obj.SetFloat(v)
		case reflect.String:
			obj.SetString(value)
		case reflect.Bool:
			var v bool
			if v, err = strconv.ParseBool(value); err != nil {
				return err
			}
			obj.SetBool(v)
		}
	case reflect.Value:
		obj.Set(value)
	default:
		obj.Set(reflect.ValueOf(value))
	}

	return nil
}
