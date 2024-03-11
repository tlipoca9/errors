package errors

import (
	"runtime"
	"strings"
)

type StackTracer interface {
	StackTrace() []StackFrame
}

var _ StackTracer = (*Error)(nil)

type StackFrame struct {
	Message string
	Frames  []runtime.Frame
}

func newStackFrame(message string, callers []uintptr) StackFrame {
	if len(callers) == 0 {
		return StackFrame{
			Message: message,
		}
	}

	runtimeFrames := runtime.CallersFrames(callers)
	frames := make([]runtime.Frame, 0, len(callers))
	for frame, more := runtimeFrames.Next(); more; frame, more = runtimeFrames.Next() {
		frames = append(frames, frame)
	}
	return StackFrame{
		Message: message,
		Frames:  frames,
	}
}

func cleanStackFrames(stackFrames []StackFrame) {
	// Remove the same frames
	for i := 1; i < len(stackFrames); i++ {
		a, b := stackFrames[i-1], stackFrames[i]
		jj := len(a.Frames) - 1
		for j, k := len(a.Frames)-1, len(b.Frames)-1; j >= 0 && k >= 0; j, k = j-1, k-1 {
			if a.Frames[j].PC != b.Frames[k].PC {
				break
			}
			jj = j - 1
		}
		// If the frames are the same, remove the frames from the previous stack
		if jj < 0 {
			stackFrames[i-1].Frames = nil
		} else {
			stackFrames[i-1].Frames = a.Frames[:jj]
		}
	}

	// Remove the runtime path
	if C.runtimePath != "" {
		for i := range stackFrames {
			for j := range stackFrames[i].Frames {
				stackFrames[i].Frames[j].File = strings.TrimPrefix(stackFrames[i].Frames[j].File, C.runtimePath)
			}
		}
	}
}
