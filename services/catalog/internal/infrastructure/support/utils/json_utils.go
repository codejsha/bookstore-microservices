package utils

import "encoding/json"

func ToJsonString(v interface{}) *string {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	s := string(b)
	return &s
}

func ParseJsonInt64Array(s *string) []int64 {
	if s == nil || *s == "" {
		return nil
	}
	var result []int64
	if err := json.Unmarshal([]byte(*s), &result); err != nil {
		return nil
	}
	return result
}

func ParseJsonStringArray(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(*s), &result); err != nil {
		return nil
	}
	return result
}
