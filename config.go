package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ErrorStyle string

const (
	StyleNormal ErrorStyle = "normal"
	StyleStack  ErrorStyle = "stack"
)

var C config

type config struct {
	runtimePath string

	Skip               int
	MaxStackDepth      uint8
	Style              ErrorStyle
	StackFramesHandler func(sfs []StackFrame) string
}

func init() {
	C = config{
		Skip:               0,
		MaxStackDepth:      64,
		Style:              StyleNormal,
		StackFramesHandler: DefaultStackFramesHandler,
	}
}

func DefaultStackFramesHandler(sfs []StackFrame) string {
	var buf bytes.Buffer
	for i, sf := range sfs {
		buf.WriteString(sf.Message)

		for _, f := range sf.Frames {
			str := fmt.Sprintf("\n  %s:%d (0x%x) %s()", f.File, f.Line, f.PC, f.Function)
			buf.WriteString(str)

		}

		if i != len(sfs)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func JSONStackFramesHandler(sfs []StackFrame) string {
	type printFrame struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	}
	type printStackFrame struct {
		Message string       `json:"message"`
		Frames  []printFrame `json:"frames,omitempty"`
	}

	ps := make([]printStackFrame, 0, len(sfs))

	for _, sf := range sfs {
		if len(sf.Frames) == 0 {
			ps = append(ps, printStackFrame{Message: sf.Message})
			continue
		}

		p := printStackFrame{
			Message: sf.Message,
			Frames:  make([]printFrame, 0, len(sf.Frames)),
		}
		for _, f := range sf.Frames {
			pf := printFrame{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			p.Frames = append(p.Frames, pf)
		}
		ps = append(ps, p)
	}

	b, _ := json.Marshal(ps)
	return string(b)
}
