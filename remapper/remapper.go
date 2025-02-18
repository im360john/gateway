package remapper

import "encoding/json"

func Remap[TValue any](config any) (TValue, error) {
	var t TValue
	raw, err := json.Marshal(config)
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return t, err
	}
	return t, nil
}
