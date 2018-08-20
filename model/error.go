package model

// ValidationError holds the validation errors of model
type ValidationError map[string][]string

func (e ValidationError) Error() string {
	return "invalid data"
}

// Add adds validation msg of the field
func (e ValidationError) Add(key, msg string) {
	e[key] = append(e[key], msg)
}
