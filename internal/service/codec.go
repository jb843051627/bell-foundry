package service

import "encoding/json"

func decode(raw []byte, target any) error { return json.Unmarshal(raw, target) }

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
