package disk

import (
	"encoding/json"
)

type Record struct {
	Index uint64 `json:"index"`
	Value string `json:"value"`
}

func (r Record) Bytes() ([]byte, error) {
	return json.Marshal(r)
}
