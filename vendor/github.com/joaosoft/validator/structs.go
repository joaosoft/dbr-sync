package validator

import (
	"reflect"

	"github.com/joaosoft/logger"
)

func (v *Validator) init() {
	v.handlersBefore = v.newDefaultBeforeHandlers()
	v.handlersMiddle = v.newDefaultMiddleHandlers()
	v.handlersAfter = v.newDefaultPosHandlers()

	v.activeHandlers = v.newActiveHandlers()
}

type Validator struct {
	tag              string
	activeHandlers   map[string]bool
	handlersBefore   map[string]beforeTagHandler
	handlersMiddle   map[string]middleTagHandler
	handlersAfter    map[string]afterTagHandler
	errorCodeHandler errorCodeHandler
	callbacks        map[string]callbackHandler
	sanitize         []string
	logger           logger.ILogger
	validateAll      bool
}

type argument struct {
	Id    string
	Value interface{}
}

func NewArgument(id string, value interface{}) *argument {
	return &argument{
		Id:    id,
		Value: value,
	}
}

type defaultValues map[string]map[string]*data

type errorCodeHandler func(context *ValidatorContext, validationData *ValidationData) error
type callbackHandler func(context *ValidatorContext, validationData *ValidationData) []error

type beforeTagHandler func(context *ValidatorContext, validationData *ValidationData) []error
type middleTagHandler func(context *ValidatorContext, validationData *ValidationData) []error
type afterTagHandler func(context *ValidatorContext, validationData *ValidationData) []error

type ValidatorContext struct {
	validator *Validator
	values    map[string]map[string]*data
}

type baseData struct {
	Id        string
	Arguments []interface{}
}

type ValidationData struct {
	*baseData
	Code           string
	Field          string
	Parent         reflect.Value
	Value          reflect.Value
	Name           string
	Expected       interface{}
	ErrorData      *errorData
	Errors         *[]error
	ErrorsReplaced map[error]bool
}

type errorData struct {
	Code      string
	Arguments []interface{}
}

type data struct {
	value reflect.Value
	typ   reflect.StructField
}

type expression struct {
	data         *data
	result       error
	expected     string
	nextOperator operator
}
