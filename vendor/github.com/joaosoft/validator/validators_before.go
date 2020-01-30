package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func (v *Validator) validate_id(context *ValidatorContext, validationData *ValidationData) []error {
	id := v._convertToString(validationData.Expected)
	dat := &data{
		value: validationData.Value,
		typ: reflect.StructField{
			Type: reflect.TypeOf(validationData.Value),
		},
	}
	context.SetValue(constTagId, id, dat)

	return nil
}

func (v *Validator) validate_if(context *ValidatorContext, validationData *ValidationData) []error {
	rtnErrs := make([]error, 0)

	str := validationData.Expected.(string)
	var expressions []*expression
	var expr *expression
	var query string

	// read conditions
	size := len(str)

	for i := 0; i < size; i++ {
		switch str[i] {
		case constParenthesesStart[0]:
			continue

		case constParenthesesEnd[0]:
			startId := strings.Index(query, fmt.Sprintf("%s=", constTagId))
			startArg := strings.Index(query, fmt.Sprintf("%s=", constTagArg))
			if startId == -1 && startArg == -1 {
				return rtnErrs
			}

			var start int
			var tag string
			if startId > -1 {
				tag = constTagId
				start = startId
			}

			if startArg > -1 {
				tag = constTagArg
				start = startArg
			}

			end := strings.Index(query[start:], " ")
			if end == -1 {
				end = size - 1
			}

			id := query[start+len(tag)+1 : end]
			query = query[end+1:]

			data, ok := context.GetValue(tag, id)

			if ok {
				var errs []error
				err := context.execute(data.typ, data.value, strings.Split(query, " "), &errs)

				// get next operator
				var oper operator
				if index := strings.Index(str[i+1:], constParenthesesStart); index > -1 {
					oper = operator(strings.TrimSpace(str[i+1 : i+1+index]))

					str = str[i+1+index:]
					i = 0
					size = len(str)
				}

				expr = &expression{
					data:         data,
					result:       err,
					nextOperator: oper,
					expected:     query,
				}
				expressions = append(expressions, expr)
			}

			query = ""

		default:
			query = fmt.Sprintf("%s%c", query, str[i])
		}
	}

	// validate all conditions
	var condition = ""
	var prevOperator = operatorNone

	for _, expr := range expressions {

		if condition == "" {
			if expr.result == nil {
				condition = constConditionOk
			} else {
				condition = constConditionKo
			}
		} else {

			switch prevOperator {
			case operatorAnd:
				if expr.result != nil {
					condition = constConditionKo
				}
			case operatorOr:
				if expr.result == nil && condition == constConditionKo {
					condition = constConditionOk
				}
			case operatorNone:
				if expr.result == nil {
					condition = constConditionOk
				}
			}
		}

		prevOperator = expr.nextOperator
	}

	if condition == constConditionKo {
		return []error{ErrorSkipValidation}
	}

	return nil
}

func (v *Validator) validate_args(context *ValidatorContext, validationData *ValidationData) []error {
	rtnErrs := make([]error, 0)

	splitArgs := strings.Split(validationData.Expected.(string), constTagSplitValues)

	for _, arg := range splitArgs {
		validationData.Arguments = append(validationData.Arguments, arg)
	}

	return rtnErrs
}
