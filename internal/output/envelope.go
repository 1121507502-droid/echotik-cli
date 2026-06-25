package output

import (
	"encoding/json"
	"io"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

type CLIError struct {
	Type    string
	Message string
	Hint    string
}

func (e *CLIError) Error() string {
	return e.Message
}

func NewError(typ, message, hint string) *CLIError {
	return &CLIError{Type: typ, Message: message, Hint: hint}
}

type Envelope struct {
	OK      bool                   `json:"ok"`
	Data    any                    `json:"data,omitempty"`
	Error   *Error                 `json:"error,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Notice  map[string]interface{} `json:"_notice,omitempty"`
	Source  string                 `json:"source,omitempty"`
	Command string                 `json:"command,omitempty"`
}

func JSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func Success(w io.Writer, data any, meta map[string]interface{}) error {
	return JSON(w, Envelope{
		OK:     true,
		Data:   data,
		Meta:   meta,
		Source: "echotik",
	})
}

func Failure(w io.Writer, typ, message, hint string) error {
	return JSON(w, Envelope{
		OK: false,
		Error: &Error{
			Type:    typ,
			Message: message,
			Hint:    hint,
		},
		Source: "echotik",
	})
}

func WriteCLIError(w io.Writer, err *CLIError) error {
	return Failure(w, err.Type, err.Message, err.Hint)
}

func ExitCode(typ string) int {
	switch typ {
	case "validation_error":
		return 2
	case "authentication_error":
		return 3
	case "rate_limit", "server_error":
		return 4
	default:
		return 1
	}
}
