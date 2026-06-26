package schema

import "strings"

func Clean(params map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range params {
		if strings.TrimSpace(v) != "" {
			out[k] = v
		}
	}
	return out
}

func RankType(raw string) string {
	switch strings.ToLower(raw) {
	case "daily", "day":
		return "1"
	case "weekly", "week":
		return "2"
	case "monthly", "month":
		return "3"
	default:
		return raw
	}
}

func PageParams(page, pageSize string) map[string]string {
	return Clean(map[string]string{
		"page_num":  page,
		"page_size": pageSize,
	})
}

func Require(name, value, example string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Message: "--" + name + " is required",
			Hint:    example,
		}
	}
	return nil
}

type ValidationError struct {
	Message string
	Hint    string
}

func (e *ValidationError) Error() string {
	return e.Message
}
