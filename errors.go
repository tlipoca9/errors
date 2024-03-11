package errors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

var _ error = (*Error)(nil)

type Error struct {
	// callers is a slice of PCs that have not yet been expanded to frames.
	callers []uintptr

	// message is the error message.
	message string

	// wrapped is the wrapped error.
	wrapped error

	// frames is the expanded stack frames.
	frames []StackFrame
}

func (e *Error) Error() string {
	switch C.Style {
	case StyleNormal:
		var buf bytes.Buffer
		buf.WriteString(e.message)

		for err := Unwrap(e); err != nil; err = Unwrap(err) {
			buf.WriteString(": ")

			if ee, ok := err.(*Error); ok {
				buf.WriteString(ee.message)
				continue
			}

			buf.WriteString(err.Error())
		}

		return buf.String()
	case StyleStack:
		return C.StackFramesHandler(e.StackTrace())
	}
	return e.message
}

func (e *Error) StackTrace() []StackFrame {
	if e.frames == nil {
		var (
			sfs     = make([]StackFrame, 0)
			collect func(err error)
		)
		collect = func(err error) {
			if err == nil {
				return
			}

			if e, ok := err.(*Error); ok {
				sfs = append(sfs, newStackFrame(e.message, e.callers))
				collect(Unwrap(err))
				return
			}

			sfs = append(sfs, newStackFrame(err.Error(), nil))

			collect(Unwrap(err))
		}

		collect(e)
		// the first frame is the top of the stack, so we reverse the slice
		for i, j := 0, len(sfs)-1; i < j; i, j = i+1, j-1 {
			sfs[i], sfs[j] = sfs[j], sfs[i]
		}
		// clean the stack frames
		cleanStackFrames(sfs)
		e.frames = sfs
	}

	return e.frames
}

func (e *Error) Unwrap() error {
	return e.wrapped
}

func New(message string) error {
	callers := make([]uintptr, C.MaxStackDepth)
	length := runtime.Callers(2+C.Skip, callers[:])
	return &Error{
		callers: callers[:length],
		message: message,
		wrapped: nil,
	}
}

func Newf(format string, a ...any) error {
	message := fmt.Sprintf(format, a...)
	callers := make([]uintptr, C.MaxStackDepth)
	length := runtime.Callers(2+C.Skip, callers[:])
	return &Error{
		callers: callers[:length],
		message: message,
		wrapped: nil,
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	callers := make([]uintptr, C.MaxStackDepth)
	length := runtime.Callers(2+C.Skip, callers[:])

	return &Error{
		callers: callers[:length],
		message: message,
		wrapped: err,
	}
}

func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, a...)
	callers := make([]uintptr, C.MaxStackDepth)
	length := runtime.Callers(2+C.Skip, callers[:])

	return &Error{
		callers: callers[:length],
		message: message,
		wrapped: err,
	}
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
