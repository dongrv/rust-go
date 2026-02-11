// package rust provides Rust-like programming constructs for Go
package rust

import (
	"fmt"
)

// Option represents an optional value: either Some(T) or None
type Option[T any] struct {
	value *T
}

// Some creates an Option containing a value
func Some[T any](value T) Option[T] {
	return Option[T]{value: &value}
}

// None creates an empty Option
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// IsSome returns true if the option is a Some value
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone returns true if the option is a None value
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Unwrap returns the contained Some value, panics if the value is None
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("called `Option.Unwrap()` on a `None` value")
	}
	return *o.value
}

// UnwrapOr returns the contained value or a provided default
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.IsSome() {
		return *o.value
	}
	return defaultValue
}

// UnwrapOrElse returns the contained value or computes it from a closure
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsSome() {
		return *o.value
	}
	return f()
}

// Expect returns the contained value or panics with a custom message
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(msg)
	}
	return *o.value
}

// Map maps an Option[T] to Option[U] by applying a function
func MapOption[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsSome() {
		return Some(f(*o.value))
	}
	return None[U]()
}

// AndThen chains operations that return Option
func AndThenOption[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsSome() {
		return f(*o.value)
	}
	return None[U]()
}

// Filter filters the Option based on a predicate
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.IsSome() && predicate(*o.value) {
		return o
	}
	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns optb
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return optb
}

// OrElse returns the option if it contains a value, otherwise calls f
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

// String returns a string representation of the Option
func (o Option[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("Some(%v)", *o.value)
	}
	return "None"
}

// Result represents a computation that may succeed (Ok) or fail (Err)
type Result[T any, E any] struct {
	ok  *T
	err *E
}

// Ok creates a successful Result containing a value
func Ok[T any, E any](value T) Result[T, E] {
	return Result[T, E]{ok: &value, err: nil}
}

// Err creates an error Result containing an error value
func Err[T any, E any](err E) Result[T, E] {
	return Result[T, E]{ok: nil, err: &err}
}

// IsOk returns true if the Result is Ok
func (r Result[T, E]) IsOk() bool {
	return r.ok != nil
}

// IsErr returns true if the Result is Err
func (r Result[T, E]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the contained Ok value, panics if the value is Err
func (r Result[T, E]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Sprintf("called `Result.Unwrap()` on an `Err` value: %v", *r.err))
	}
	return *r.ok
}

// UnwrapErr returns the contained Err value, panics if the value is Ok
func (r Result[T, E]) UnwrapErr() E {
	if r.IsOk() {
		panic(fmt.Sprintf("called `Result.UnwrapErr()` on an `Ok` value: %v", *r.ok))
	}
	return *r.err
}

// UnwrapOr returns the contained Ok value or a provided default
func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if r.IsOk() {
		return *r.ok
	}
	return defaultValue
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure
func (r Result[T, E]) UnwrapOrElse(f func(E) T) T {
	if r.IsOk() {
		return *r.ok
	}
	return f(*r.err)
}

// Expect returns the contained Ok value or panics with a custom message
func (r Result[T, E]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Sprintf("%s: %v", msg, *r.err))
	}
	return *r.ok
}

// Map maps a Result[T, E] to Result[U, E] by applying a function
func MapResult[T any, E any, U any](r Result[T, E], f func(T) U) Result[U, E] {
	if r.IsOk() {
		return Ok[U, E](f(*r.ok))
	}
	return Err[U, E](*r.err)
}

// MapErr maps a Result[T, E] to Result[T, F] by applying a function to error
func MapErrResult[T any, E any, F any](r Result[T, E], f func(E) F) Result[T, F] {
	if r.IsErr() {
		return Err[T, F](f(*r.err))
	}
	return Ok[T, F](*r.ok)
}

// AndThen chains operations that return Result
func AndThenResult[T any, E any, U any](r Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	if r.IsOk() {
		return f(*r.ok)
	}
	return Err[U, E](*r.err)
}

// Or returns the result if it contains an Ok value, otherwise returns resb
func (r Result[T, E]) Or(resb Result[T, E]) Result[T, E] {
	if r.IsOk() {
		return r
	}
	return resb
}

// OrElse returns the result if it contains an Ok value, otherwise calls f
func (r Result[T, E]) OrElse(f func(E) Result[T, E]) Result[T, E] {
	if r.IsOk() {
		return r
	}
	return f(*r.err)
}

// String returns a string representation of the Result
func (r Result[T, E]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok(%v)", *r.ok)
	}
	return fmt.Sprintf("Err(%v)", *r.err)
}
