// Package errors provides enhanced error handling utilities inspired by Rust's error handling patterns.
// This package reduces boilerplate code by providing functional error handling patterns,
// error composition, and context propagation.
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Error is an enhanced error type that supports chaining, context, and structured error information.
type Error struct {
	// Message is the human-readable error message
	Message string

	// Cause is the underlying error that caused this error
	Cause error

	// Stack contains the call stack when the error was created
	Stack []uintptr

	// Context contains additional structured context about the error
	Context map[string]interface{}
}

// New creates a new error with the given message.
func New(message string) *Error {
	return &Error{
		Message: message,
		Stack:   captureStack(2), // Skip New and the caller
		Context: make(map[string]interface{}),
	}
}

// Errorf creates a new error with formatted message.
func Errorf(format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Stack:   captureStack(2),
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context.
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Message: message + ": " + err.Error(),
		Cause:   err,
		Stack:   captureStack(2),
		Context: make(map[string]interface{}),
	}
}

// Wrapf wraps an existing error with formatted message.
func Wrapf(err error, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)
	return &Error{
		Message: message + ": " + err.Error(),
		Cause:   err,
		Stack:   captureStack(2),
		Context: make(map[string]interface{}),
	}
}

// WithContext adds structured context to the error.
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithContextMap adds multiple context key-value pairs to the error.
func (e *Error) WithContextMap(context map[string]interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	for k, v := range context {
		e.Context[k] = v
	}
	return e
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Cause
}

// StackTrace returns the stack trace as a formatted string.
func (e *Error) StackTrace() string {
	if len(e.Stack) == 0 {
		return ""
	}

	frames := runtime.CallersFrames(e.Stack)
	var sb strings.Builder
	sb.WriteString("Stack trace:\n")

	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("  %s\n    %s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}

	return sb.String()
}

// String returns a detailed string representation of the error.
func (e *Error) String() string {
	var sb strings.Builder
	sb.WriteString(e.Error())

	if len(e.Context) > 0 {
		sb.WriteString("\nContext:")
		for k, v := range e.Context {
			sb.WriteString(fmt.Sprintf("\n  %s: %v", k, v))
		}
	}

	if len(e.Stack) > 0 {
		sb.WriteString("\n")
		sb.WriteString(e.StackTrace())
	}

	return sb.String()
}

// Result is a type alias for functions that return a value and an error.
// It enables functional error handling patterns.
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Map applies a function to the value if the Result is Ok.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(f(r.value))
}

// MapErr applies a function to transform the error if the Result is Err.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return Err[T](f(r.err))
}

// AndThen chains operations that return Result.
func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return f(r.value)
}

// OrElse returns the Result if it's Ok, otherwise calls f.
func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if r.err == nil {
		return r
	}
	return f(r.err)
}

// Unwrap returns the value or panics if there's an error.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("called Result.Unwrap() on error: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the value or a default if there's an error.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value or computes a default from the error.
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.err != nil {
		return f(r.err)
	}
	return r.value
}

// Expect returns the value or panics with a custom message if there's an error.
func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}
	return r.value
}

// IsOk returns true if the Result is Ok.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result is Err.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Error returns the error if the Result is Err.
func (r Result[T]) Error() error {
	return r.err
}

// Value returns the value and error separately.
func (r Result[T]) Value() (T, error) {
	return r.value, r.err
}

// Try is a helper function that converts a function returning (T, error) to Result[T].
func Try[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// TryFunc wraps a function that returns (T, error) and returns a Result[T].
func TryFunc[T any](f func() (T, error)) Result[T] {
	value, err := f()
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Recover converts a panic to an error Result.
func Recover[T any](f func() T) (result Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			result = Err[T](fmt.Errorf("panic recovered: %v", r))
		}
	}()

	return Ok(f())
}

// Combine combines multiple Results into a single Result of slice.
func Combine[T any](results ...Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if r.err != nil {
			return Err[[]T](r.err)
		}
		values = append(values, r.value)
	}
	return Ok(values)
}

// FirstError returns the first error from multiple Results.
func FirstError[T any](results ...Result[T]) error {
	for _, r := range results {
		if r.err != nil {
			return r.err
		}
	}
	return nil
}

// AllOk returns true if all Results are Ok.
func AllOk[T any](results ...Result[T]) bool {
	for _, r := range results {
		if r.err != nil {
			return false
		}
	}
	return true
}

// AnyErr returns true if any Result is Err.
func AnyErr[T any](results ...Result[T]) bool {
	for _, r := range results {
		if r.err != nil {
			return true
		}
	}
	return false
}

// captureStack captures the current call stack.
func captureStack(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[:n]
}

// ErrorHandler provides a fluent interface for error handling.
type ErrorHandler struct {
	err  error
	skip bool
}

// Handle creates a new ErrorHandler for the given error.
func Handle(err error) *ErrorHandler {
	return &ErrorHandler{err: err}
}

// If returns the ErrorHandler if the condition is true.
func (h *ErrorHandler) If(condition bool) *ErrorHandler {
	if !condition {
		// Return a new handler that will skip subsequent Then calls
		return &ErrorHandler{err: h.err, skip: true}
	}
	return h
}

// Then executes the function if there's no error and not skipping.
func (h *ErrorHandler) Then(f func() error) *ErrorHandler {
	if h.err != nil || h.skip {
		return h
	}
	h.err = f()
	return h
}

// Map transforms the error using the provided function.
func (h *ErrorHandler) Map(f func(error) error) *ErrorHandler {
	if h.err != nil && !h.skip {
		h.err = f(h.err)
	}
	return h
}

// OrElse returns the error or executes the function if there's no error.
func (h *ErrorHandler) OrElse(f func() error) error {
	if h.err != nil || h.skip {
		return h.err
	}
	return f()
}

// Unwrap returns the error.
func (h *ErrorHandler) Unwrap() error {
	return h.err
}

// Must panics if there's an error.
func (h *ErrorHandler) Must() {
	if h.err != nil && !h.skip {
		panic(h.err)
	}
}

// MustWith panics with custom message if there's an error.
func (h *ErrorHandler) MustWith(msg string) {
	if h.err != nil && !h.skip {
		panic(fmt.Sprintf("%s: %v", msg, h.err))
	}
}

// Ignore ignores the error.
func (h *ErrorHandler) Ignore() {
	// Do nothing - error is ignored
}

// Log logs the error if it exists.
func (h *ErrorHandler) Log(logger func(string, ...interface{})) {
	if h.err != nil && !h.skip {
		logger("Error: %v", h.err)
	}
}

// ErrorChain represents a chain of errors for detailed error tracing.
type ErrorChain []error

// NewChain creates a new ErrorChain.
func NewChain(errors ...error) ErrorChain {
	return errors
}

// Append appends an error to the chain.
func (c ErrorChain) Append(err error) ErrorChain {
	return append(c, err)
}

// Last returns the last error in the chain.
func (c ErrorChain) Last() error {
	if len(c) == 0 {
		return nil
	}
	return c[len(c)-1]
}

// First returns the first error in the chain.
func (c ErrorChain) First() error {
	if len(c) == 0 {
		return nil
	}
	return c[0]
}

// String returns a string representation of the error chain.
func (c ErrorChain) String() string {
	if len(c) == 0 {
		return "no errors"
	}

	var sb strings.Builder
	sb.WriteString("Error chain:\n")
	for i, err := range c {
		sb.WriteString(fmt.Sprintf("  [%d] %v\n", i, err))
	}
	return sb.String()
}

// Flatten flattens nested errors into a single chain.
func Flatten(err error) ErrorChain {
	var chain ErrorChain
	for err != nil {
		chain = append(chain, err)
		if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
			err = unwrapper.Unwrap()
		} else {
			break
		}
	}
	return chain
}
