package json

func Marshal(object interface{}, tags ...string) ([]byte, error) {
	return newMarshal(object, tags...).execute()
}

func Unmarshal(bytes []byte, object interface{}, tags ...string) error {
	return newUnmarshal(bytes, object, tags...).execute()
}
