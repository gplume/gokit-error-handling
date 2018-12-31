package errs

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Error ...
type Error struct {
	Err     error
	Message string
	Code    int
	Caller  string
	Stack   *Stack
}

// assert Error implements the error interface.
// so you can pass *Error to the normal 'error' of std library
// waited by your usual code or libraries like go-kit...
var _ error = &Error{}

// Error implements the error interface.
func (e *Error) Error() string {
	b := new(bytes.Buffer)
	e.printStack(b)
	pad(b, ": ")
	b.WriteString(e.Message)
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// New intanciates the error on place so we have a precise caller
// note the now *Error pass through as an std 'error' type...
func New(err error, msg string, code int) error {
	_, w, ln, ok := runtime.Caller(1)
	er := &Error{
		Code:    code,
		Message: msg,
		Err:     err,
	}
	if ok {
		er.Caller = fmt.Sprintf("%s:%d", path.Base(w), ln)
	}
	er.populateStack()
	return er

}

// Stack ...
type Stack struct {
	Callers []uintptr
}

// populateStack uses the runtime to populate the Error's stack struct with
// information about the current stack.
func (e *Error) populateStack() {
	e.Stack = &Stack{Callers: callers()}
}

// printStack formats and prints the stack for this Error to the given buffer.
// It should be called from the Error's Error method.
func (e *Error) printStack(b *bytes.Buffer) {
	if e.Stack == nil {
		return
	}

	printCallers := callers()

	// Iterate backward through e.Stack.Callers (the last in the stack is the
	// earliest call, such as main) skipping over the PCs that are shared
	// by the error stack and by this function call stack, printing the
	// names of the functions and their file names and line numbers.
	var prev string // the name of the last-seen function
	var diff bool   // do the print and error call stacks differ now?
	for i := 0; i < len(e.Stack.Callers); i++ {
		thisFrame := frame(e.Stack.Callers, i)
		name := thisFrame.Func.Name()

		if !diff && i < len(printCallers) {
			if name == frame(printCallers, i).Func.Name() {
				// both stacks share this PC, skip it.
				continue
			}
			// No match, don't consider printCallers again.
			diff = true
		}

		// Don't print the same function twice.
		// (Can happen when multiple error stacks have been coalesced.)
		if name == prev {
			continue
		}

		// Find the uncommon prefix between this and the previous
		// function name, separating by dots and slashes.
		trim := 0
		for {
			j := strings.IndexAny(name[trim:], "./")
			if j < 0 {
				break
			}
			if !strings.HasPrefix(prev, name[:j+trim]) {
				break
			}
			trim += j + 1 // skip over the separator
		}

		// Do the printing.
		pad(b, separator)
		fmt.Fprintf(b, "%v:%d: ", thisFrame.File, thisFrame.Line)
		if trim > 0 {
			b.WriteString("...")
		}
		b.WriteString(name[trim:])

		prev = name
	}
}

// frame returns the nth frame, with the frame at top of stack being 0.
func frame(callers []uintptr, n int) *runtime.Frame {
	frames := runtime.CallersFrames(callers)
	var f runtime.Frame
	for i := len(callers) - 1; i >= n; i-- {
		var ok bool
		f, ok = frames.Next()
		if !ok {
			break // Should never happen, and this is just debugging.
		}
	}
	return &f
}

// callers is a wrapper for runtime.callers that allocates a slice.
func callers() []uintptr {
	var stk [64]uintptr
	const skip = 4 // Skip 4 stack frames; ok for both E and Error funcs.
	n := runtime.Callers(skip, stk[:])
	return stk[:n]
}

var separator = ":\n\t"

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
