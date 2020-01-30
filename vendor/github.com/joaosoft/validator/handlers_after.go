package validator

func (v *Validator) newDefaultPosHandlers() map[string]afterTagHandler {
	return map[string]afterTagHandler{
		constTagError: v.validate_error,
	}
}
