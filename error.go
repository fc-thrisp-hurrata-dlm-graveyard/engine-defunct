package engine

import (
	"bytes"
	"fmt"
)

const (
	ErrorTypeInternal = 1 << iota
	ErrorTypeExternal = 1 << iota
	ErrorTypeAll      = 0xffffffff
)

type (
	// Used with Ctx to collect errors that occurred during a http request.
	errorMsg struct {
		Err  string      `json:"error"`
		Type uint32      `json:"-"`
		Meta interface{} `json:"meta"`
	}

	errorMsgs []errorMsg
)

func (a errorMsgs) ByType(typ uint32) errorMsgs {
	if len(a) == 0 {
		return a
	}
	result := make(errorMsgs, 0, len(a))
	for _, msg := range a {
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}
	return result
}

func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		text := fmt.Sprintf("Engine Error #%02d: %s \n     Meta: %v\n", (i + 1), msg.Err, msg.Meta)
		buffer.WriteString(text)
	}
	return buffer.String()
}
