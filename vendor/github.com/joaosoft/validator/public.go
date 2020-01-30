package validator

func AddBefore(name string, handler beforeTagHandler) *Validator {
	return validatorInstance.AddBefore(name, handler)
}

func AddMiddle(name string, handler middleTagHandler) *Validator {
	return validatorInstance.AddMiddle(name, handler)
}

func AddAfter(name string, handler afterTagHandler) *Validator {
	return validatorInstance.AddAfter(name, handler)
}

func SetErrorCodeHandler(handler errorCodeHandler) *Validator {
	return validatorInstance.SetErrorCodeHandler(handler)
}

func SetValidateAll(validate bool) *Validator {
	return validatorInstance.SetValidateAll(validate)
}

func SetTag(tag string) *Validator {
	return validatorInstance.SetTag(tag)
}

func SetSanitize(sanitize []string) *Validator {
	return validatorInstance.SetSanitize(sanitize)
}

func AddCallback(name string, callback callbackHandler) *Validator {
	return validatorInstance.AddCallback(name, callback)
}

func Validate(obj interface{}, args ...*argument) []error {
	return NewValidatorHandler(validatorInstance, args...).handleValidation(obj)
}
