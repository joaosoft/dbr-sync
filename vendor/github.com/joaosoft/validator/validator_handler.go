package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func (v *Validator) newDefaultValues() defaultValues {
	return map[string]map[string]*data{
		constTagId:   make(map[string]*data),
		constTagJson: make(map[string]*data),
		constTagArg:  make(map[string]*data),
	}
}

func NewValidatorHandler(validator *Validator, args ...*argument) *ValidatorContext {
	context := &ValidatorContext{
		validator: validator,
		values:    validator.newDefaultValues(),
	}

	for _, arg := range args {
		context.values[constTagArg][arg.Id] = &data{
			value: reflect.ValueOf(arg.Value),
			typ: reflect.StructField{
				Type: reflect.TypeOf(arg.Value),
			},
		}
	}

	return context
}

func (vc *ValidatorContext) GetValue(tag string, id string) (*data, bool) {
	if values, ok := vc.values[tag]; ok {
		if value, ok := values[id]; ok {
			return value, ok
		}
	}
	return nil, false
}

func (vc *ValidatorContext) SetValue(tag string, id string, value *data) bool {
	if values, ok := vc.values[tag]; ok {
		values[id] = value
		return true
	}
	return false
}

func (vc *ValidatorContext) handleValidation(value interface{}) []error {
	var err error
	errs := make([]error, 0)

	// execute
	if err = vc.do(reflect.ValueOf(value), &errs); err != nil {
		return []error{err}
	}

	return errs
}

func (vc *ValidatorContext) _getValue(value reflect.Value) (reflect.Type, reflect.Value, error) {
	types := reflect.TypeOf(value.Interface())

again:
	if (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) && !value.IsNil() {
		value = value.Elem()
		types = value.Type()
		goto again
	}

	return types, value, nil
}

func (vc *ValidatorContext) load(value reflect.Value, errs *[]error) (err error) {
	var types reflect.Type
	types, value, err = vc._getValue(value)
	if err != nil {
		return err
	}

	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < types.NumField(); i++ {
			var dat *data
			nextValue := value.Field(i)
			nextType := types.Field(i)

			if !nextValue.CanInterface() {
				continue
			}

			tagValue, exists := nextType.Tag.Lookup(vc.validator.tag)

			// save id sub tags
			var err error
			if exists && strings.Contains(tagValue, fmt.Sprintf("%s=", constTagId)) {
				var id string

				split := strings.Split(tagValue, ",")
				var tag []string
				for _, item := range split {
					tag = strings.Split(item, "=")
					tag[0] = strings.TrimSpace(tag[0])

					switch tag[0] {
					case constTagId:
						id = tag[1]
						if dat == nil {
							dat = &data{
								value: nextValue,
								typ:   nextType,
							}
						}
					case constTagSet:
						newStruct := reflect.New(value.Type()).Elem()
						newField := newStruct.Field(i)

						if !strings.Contains(tagValue, fmt.Sprintf("%s=", constTagIf)) {
							if err = _setValue(nextValue.Kind(), newField, tag[1]); err != nil {
								*errs = append(*errs, err)
							}
						} else {
							if err = _setValue(nextValue.Kind(), newField, value.Field(i)); err != nil {
								*errs = append(*errs, err)
							}
						}

						if len(*errs) > 0 && !vc.validator.validateAll {
							return nil
						}

						dat = &data{
							value: newField,
							typ:   nextType,
						}
					}
				}
				vc.SetValue(tag[0], id, dat)
			}

			// save json tags
			tagValue, exists = nextType.Tag.Lookup(constTagJson)
			if exists && tagValue != "-" {
				split := strings.Split(tagValue, ",")
				dat = &data{
					value: nextValue,
					typ:   nextType,
				}
				vc.SetValue(constTagJson, split[0], dat)
			}

			if err := vc.load(nextValue, errs); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := vc.load(nextValue, errs); err != nil {
				return err
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if err := vc.load(key, errs); err != nil {
				return err
			}
			if err := vc.load(nextValue, errs); err != nil {
				return err
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func (vc *ValidatorContext) do(value reflect.Value, errs *[]error) (err error) {
	var types reflect.Type
	types, value, err = vc._getValue(value)
	if err != nil {
		return err
	}

	switch value.Kind() {
	case reflect.Struct:

		// load id's
		if err := vc.load(reflect.ValueOf(value), errs); err != nil {
			return err
		}

		for i := 0; i < types.NumField(); i++ {
			nextValue := value.Field(i)
			nextType := types.Field(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := vc.doValidate(nextValue, nextType, errs); err != nil {
				return err
			}

			if len(*errs) > 0 && !vc.validator.validateAll {
				return nil
			}

			if err := vc.do(nextValue, errs); err != nil {
				return err
			}

			if len(*errs) > 0 && !vc.validator.validateAll {
				return nil
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)

			if !nextValue.CanInterface() {
				continue
			}

			if err := vc.do(nextValue, errs); err != nil {
				return err
			}

			if len(*errs) > 0 && !vc.validator.validateAll {
				return nil
			}
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			nextValue := value.MapIndex(key)

			if !nextValue.CanInterface() {
				continue
			}

			if err := vc.do(key, errs); err != nil {
				return err
			}

			if len(*errs) > 0 && !vc.validator.validateAll {
				return nil
			}

			if err := vc.do(nextValue, errs); err != nil {
				return err
			}

			if len(*errs) > 0 && !vc.validator.validateAll {
				return nil
			}
		}

	default:
		// do nothing ...
	}
	return nil
}

func (vc *ValidatorContext) doValidate(value reflect.Value, typ reflect.StructField, errs *[]error) error {

	tag, exists := typ.Tag.Lookup(vc.validator.tag)
	if !exists {
		return nil
	}

	validations := strings.Split(tag, ",")

	return vc.execute(typ, value, validations, errs)
}

func (vc *ValidatorContext) getFieldId(validations []string) string {
	for _, validation := range validations {
		options := strings.SplitN(validation, "=", 2)
		tag := strings.TrimSpace(options[0])

		if tag == constTagId {
			return options[1]
		}
	}

	return ""
}

func (vc *ValidatorContext) execute(typ reflect.StructField, value reflect.Value, validations []string, errs *[]error) error {
	var err error
	var itErrs []error
	var replacedErrors = make(map[error]bool)
	skipValidation := false
	onlyHandleNextErrorTag := false

	defer func() {
		*errs = append(*errs, itErrs...)
	}()

	baseData := &baseData{
		Id:        vc.getFieldId(validations),
		Arguments: make([]interface{}, 0),
	}

	for _, validation := range validations {
		var name string
		var tag string
		var prefix string

		options := strings.SplitN(validation, "=", 2)
		tag = strings.TrimSpace(options[0])

		if split := strings.SplitN(tag, ":", 2); len(split) > 1 {
			prefix = split[0]
			tag = split[1]
		}

		if onlyHandleNextErrorTag && !vc.validator.validateAll && tag != constTagError {
			continue
		}

		if _, ok := vc.validator.activeHandlers[tag]; !ok {
			return fmt.Errorf("invalid tag [%s]", tag)
		}

		var expected interface{}
		if len(options) > 1 {
			expected = strings.TrimSpace(options[1])
		}

		jsonName, exists := typ.Tag.Lookup(constTagJson)
		if exists {
			split := strings.SplitN(jsonName, ",", 2)
			name = split[0]
		} else {
			name = typ.Name
		}

		if skipValidation {
			if tag == constTagIf {
				skipValidation = false
			} else {
				continue
			}
		}

		// execute validations
		switch prefix {
		case constPrefixTagKey, constPrefixTagItem:
			if !value.CanInterface() {
				return nil
			}

			var types reflect.Type
			types, value, err = vc._getValue(value)
			if err != nil {
				return err
			}

			if prefix == constPrefixTagKey && value.Kind() != reflect.Map {
				continue
			}

			switch value.Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < value.Len(); i++ {
					nextValue := value.Index(i)

					if !nextValue.CanInterface() {
						continue
					}

					validationData := ValidationData{
						baseData:       baseData,
						Name:           name,
						Field:          typ.Name,
						Parent:         value,
						Value:          nextValue,
						Expected:       expected,
						Errors:         &itErrs,
						ErrorsReplaced: replacedErrors,
					}

					err = vc.executeHandlers(tag, &validationData, &itErrs)
				}
			case reflect.Map:
				for _, key := range value.MapKeys() {

					var nextValue reflect.Value

					switch prefix {
					case constPrefixTagKey:
						nextValue = key
					case constPrefixTagItem:
						nextValue = value.MapIndex(key)
					}

					if !key.CanInterface() {
						continue
					}

					validationData := ValidationData{
						baseData:       baseData,
						Name:           name,
						Field:          typ.Name,
						Parent:         value,
						Value:          nextValue,
						Expected:       expected,
						Errors:         &itErrs,
						ErrorsReplaced: replacedErrors,
					}

					err = vc.executeHandlers(tag, &validationData, &itErrs)
				}
			case reflect.Struct:
				for i := 0; i < types.NumField(); i++ {
					nextValue := value.Field(i)

					if !nextValue.CanInterface() {
						continue
					}

					validationData := ValidationData{
						baseData:       baseData,
						Name:           name,
						Field:          typ.Name,
						Parent:         value,
						Value:          nextValue,
						Expected:       expected,
						Errors:         &itErrs,
						ErrorsReplaced: replacedErrors,
					}

					err = vc.executeHandlers(tag, &validationData, &itErrs)
				}
			}

		default:
			if prefix != "" {
				return fmt.Errorf("invalid tag prefix [%s] on tag [%s]", prefix, tag)
			}

			validationData := ValidationData{
				baseData:       baseData,
				Name:           name,
				Field:          typ.Name,
				Parent:         value,
				Value:          value,
				Expected:       expected,
				Errors:         &itErrs,
				ErrorsReplaced: replacedErrors,
			}

			err = vc.executeHandlers(tag, &validationData, &itErrs)
		}

		if onlyHandleNextErrorTag && !vc.validator.validateAll && tag == constTagError {
			if err == ErrorSkipValidation {
				skipValidation = true
				continue
			}

			return nil
		}

		if err != nil {
			if err == ErrorSkipValidation {
				skipValidation = true
				continue
			} else {
				return err
			}
		}

		if len(*errs) > 0 {
			if !onlyHandleNextErrorTag && !vc.validator.validateAll && tag != constTagError {
				onlyHandleNextErrorTag = true
				continue
			}

			if !vc.validator.validateAll {
				return nil
			}
		}
	}

	return nil
}

func (vc *ValidatorContext) executeHandlers(tag string, validationData *ValidationData, errs *[]error) error {
	var err error

	if _, ok := vc.validator.handlersBefore[tag]; ok {
		if rtnErrs := vc.validator.handlersBefore[tag](vc, validationData); rtnErrs != nil && len(rtnErrs) > 0 {

			// skip validation
			if rtnErrs[0] == ErrorSkipValidation {
				return rtnErrs[0]
			}
			*errs = append(*errs, rtnErrs...)
		}
	}

	if _, ok := vc.validator.handlersMiddle[tag]; ok {
		if rtnErrs := vc.validator.handlersMiddle[tag](vc, validationData); rtnErrs != nil && len(rtnErrs) > 0 {
			*errs = append(*errs, rtnErrs...)
		}
	}

	if _, ok := vc.validator.handlersAfter[tag]; ok {
		if rtnErrs := vc.validator.handlersAfter[tag](vc, validationData); rtnErrs != nil && len(rtnErrs) > 0 {
			*errs = append(*errs, rtnErrs...)
		}
	}

	return err
}
