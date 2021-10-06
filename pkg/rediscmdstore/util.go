package rediscmdstore

import (
	"encoding/json"
)

func mapToString(m map[string]string) (res string, err error) {
	b, err := json.Marshal(m)
	if err == nil {
		res = string(b)
	}
	return
}

func stringToMap(v string) (m map[string]string, err error) {
	m = make(map[string]string)
	err = json.Unmarshal([]byte(v), &m)
	return
}
