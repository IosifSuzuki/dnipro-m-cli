package jsonx

import "encoding/json"

type NullableFloat64 struct {
	Value *float64
}

func (n *NullableFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `"null"` {
		n.Value = nil
		return nil
	}

	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.Value = &v
	return nil
}
