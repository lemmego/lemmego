package shared

import "encoding/json"

type ValidationErrors map[string][]string

func (e ValidationErrors) Error() string {
	val, _ := json.Marshal(e)
	return string(val)
}
