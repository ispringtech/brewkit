package maybe

import (
	"encoding/json"
)

func (m *Maybe[T]) UnmarshalJSON(bytes []byte) error {
	var v T
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}

	m.v = v
	m.valid = true

	return nil
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if !m.valid {
		return []byte("null"), nil
	}
	return json.Marshal(m.v)
}
