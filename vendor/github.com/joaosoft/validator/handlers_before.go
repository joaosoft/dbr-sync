package validator

func (v *Validator) newDefaultBeforeHandlers() map[string]beforeTagHandler {
	return map[string]beforeTagHandler{
		constTagId:   v.validate_id,
		constTagIf:   v.validate_if,
		constTagArgs: v.validate_args,
	}
}
