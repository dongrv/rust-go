// Package pattern provides Rust-like pattern matching for Go
// This library brings expressive pattern matching capabilities to Go,
// allowing developers to write more declarative and readable code.
//
// # Examples
//
// Basic usage:
//
//	value := Some(42)
//	Match(value).
//		Some(func(x int) {
//			fmt.Printf("Got value: %d\n", x)
//		}).
//		None(func() {
//			fmt.Println("Got nothing")
//		})
//
// Pattern matching with Result:
//
//	result := Ok[int, string](42)
//	Match(result).
//		Ok(func(x int) {
//			fmt.Printf("Success: %d\n", x)
//		}).
//		Err(func(err string) {
//			fmt.Printf("Error: %s\n", err)
//		})
//
// Exhaustiveness checking:
//
//	Match(value).
//		Some(func(x int) { ... }).
//		None(func() { ... }).
//		Exhaustive() // Ensures all cases are handled
package pattern

import (
	"fmt"
	"reflect"
)

// Matcher is the main type for pattern matching.
// It holds the value to match against and tracks whether a match has been made.
type Matcher struct {
	value   interface{}
	matched bool
}

// Match creates a new Matcher for the given value.
// This is the entry point for pattern matching.
//
// Example:
//
//	Match(Some(42)).
//		Some(func(x int) { ... }).
//		None(func() { ... })
func Match(value interface{}) *Matcher {
	return &Matcher{
		value:   value,
		matched: false,
	}
}

// Some matches an Option[T] that contains a value.
// It executes the provided function if the Option is Some.
//
// Example:
//
//	Match(Some(42)).
//		Some(func(x int) {
//			fmt.Printf("Got: %d\n", x)
//		})
func (m *Matcher) Some(f interface{}) *Matcher {
	if m.matched {
		return m
	}

	val := reflect.ValueOf(m.value)
	// Check if it's an Option type by looking for IsSome method
	isSomeMethod := val.MethodByName("IsSome")
	if isSomeMethod.IsValid() {
		results := isSomeMethod.Call(nil)
		if len(results) > 0 && results[0].Bool() {
			// Get the value from Some
			unwrapMethod := val.MethodByName("Unwrap")
			if unwrapMethod.IsValid() {
				results := unwrapMethod.Call(nil)
				if len(results) > 0 {
					// Call the provided function with the unwrapped value
					fv := reflect.ValueOf(f)
					if fv.Kind() == reflect.Func {
						fv.Call([]reflect.Value{results[0]})
					}
					m.matched = true
				}
			}
		}
	}
	return m
}

// None matches an Option[T] that is empty.
// It executes the provided function if the Option is None.
//
// Example:
//
//	Match(None[int]()).
//		None(func() {
//			fmt.Println("Got nothing")
//		})
func (m *Matcher) None(f func()) *Matcher {
	if m.matched {
		return m
	}

	val := reflect.ValueOf(m.value)
	isNoneMethod := val.MethodByName("IsNone")
	if isNoneMethod.IsValid() {
		results := isNoneMethod.Call(nil)
		if len(results) > 0 && results[0].Bool() {
			f()
			m.matched = true
		}
	}
	return m
}

// Ok matches a Result[T, E] that contains a success value.
// It executes the provided function if the Result is Ok.
//
// Example:
//
//	Match(Ok[int, string](42)).
//		Ok(func(x int) {
//			fmt.Printf("Success: %d\n", x)
//		})
func (m *Matcher) Ok(f interface{}) *Matcher {
	if m.matched {
		return m
	}

	val := reflect.ValueOf(m.value)
	isOkMethod := val.MethodByName("IsOk")
	if isOkMethod.IsValid() {
		results := isOkMethod.Call(nil)
		if len(results) > 0 && results[0].Bool() {
			unwrapMethod := val.MethodByName("Unwrap")
			if unwrapMethod.IsValid() {
				results := unwrapMethod.Call(nil)
				if len(results) > 0 {
					// Call the provided function with the unwrapped value
					fv := reflect.ValueOf(f)
					if fv.Kind() == reflect.Func {
						fv.Call([]reflect.Value{results[0]})
					}
					m.matched = true
				}
			}
		}
	}
	return m
}

// Err matches a Result[T, E] that contains an error.
// It executes the provided function if the Result is Err.
//
// Example:
//
//	Match(Err[int, string]("error")).
//		Err(func(err string) {
//			fmt.Printf("Error: %s\n", err)
//		})
func (m *Matcher) Err(f interface{}) *Matcher {
	if m.matched {
		return m
	}

	val := reflect.ValueOf(m.value)
	isErrMethod := val.MethodByName("IsErr")
	if isErrMethod.IsValid() {
		results := isErrMethod.Call(nil)
		if len(results) > 0 && results[0].Bool() {
			unwrapErrMethod := val.MethodByName("UnwrapErr")
			if unwrapErrMethod.IsValid() {
				results := unwrapErrMethod.Call(nil)
				if len(results) > 0 {
					// Call the provided function with the error value
					fv := reflect.ValueOf(f)
					if fv.Kind() == reflect.Func {
						fv.Call([]reflect.Value{results[0]})
					}
					m.matched = true
				}
			}
		}
	}
	return m
}

// Value matches any value that equals the expected value.
// It executes the provided function if the values are equal.
//
// Example:
//
//	Match(42).
//		Value(42, func() {
//			fmt.Println("Got exactly 42")
//		})
func (m *Matcher) Value(expected interface{}, f func()) *Matcher {
	if m.matched {
		return m
	}

	if reflect.DeepEqual(m.value, expected) {
		f()
		m.matched = true
	}
	return m
}

// Type matches based on the type of the value.
// It executes the provided function if the value can be converted to the target type.
//
// Example:
//
//	var value interface{} = "hello"
//	Match(value).
//		Type(func(s string) {
//			fmt.Printf("String: %s\n", s)
//		})
func (m *Matcher) Type(f interface{}) *Matcher {
	if m.matched {
		return m
	}

	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return m
	}

	// Get the function's first parameter type
	ft := fv.Type()
	if ft.NumIn() != 1 {
		return m
	}

	targetType := ft.In(0)
	val := reflect.ValueOf(m.value)

	if val.Type().ConvertibleTo(targetType) {
		converted := val.Convert(targetType)
		fv.Call([]reflect.Value{converted})
		m.matched = true
	}
	return m
}

// Predicate matches based on a custom predicate function.
// It executes the provided function if the predicate returns true.
//
// Example:
//
//	Match(42).
//		Predicate(func(x int) bool { return x%2 == 0 }, func() {
//			fmt.Println("Even number")
//		})
func (m *Matcher) Predicate(pred interface{}, f func()) *Matcher {
	if m.matched {
		return m
	}

	pv := reflect.ValueOf(pred)
	if pv.Kind() != reflect.Func {
		return m
	}

	// Check if predicate accepts the value type
	pt := pv.Type()
	if pt.NumIn() != 1 {
		return m
	}

	val := reflect.ValueOf(m.value)
	if val.Type().ConvertibleTo(pt.In(0)) {
		converted := val.Convert(pt.In(0))
		results := pv.Call([]reflect.Value{converted})
		if len(results) > 0 && results[0].Bool() {
			f()
			m.matched = true
		}
	}
	return m
}

// Default provides a fallback case when no other patterns match.
// It should always be the last case in a match expression.
//
// Example:
//
//	Match(value).
//		Some(func(x int) { ... }).
//		None(func() { ... }).
//		Default(func() {
//			fmt.Println("Fallback")
//		})
func (m *Matcher) Default(f func()) *Matcher {
	if !m.matched {
		f()
		m.matched = true
	}
	return m
}

// Exhaustive ensures that all possible cases have been handled.
// It panics if no match was made.
//
// Example:
//
//	Match(value).
//		Some(func(x int) { ... }).
//		None(func() { ... }).
//		Exhaustive()
func (m *Matcher) Exhaustive() {
	if !m.matched {
		panic(fmt.Sprintf("pattern: non-exhaustive match on value: %v", m.value))
	}
}

// Map transforms the value using the provided function.
// Returns a new Matcher with the transformed value.
//
// Example:
//
//	result := Match(42).
//		Map(func(x int) string {
//			return fmt.Sprintf("Number: %d", x)
//		}).
//		UnwrapOr("default")
func (m *Matcher) Map(f interface{}) *Matcher {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return &Matcher{value: nil, matched: m.matched}
	}

	ft := fv.Type()
	if ft.NumIn() != 1 {
		return &Matcher{value: nil, matched: m.matched}
	}

	val := reflect.ValueOf(m.value)
	if val.Type().ConvertibleTo(ft.In(0)) {
		converted := val.Convert(ft.In(0))
		results := fv.Call([]reflect.Value{converted})
		if len(results) > 0 {
			return &Matcher{value: results[0].Interface(), matched: m.matched}
		}
	}

	return &Matcher{value: nil, matched: m.matched}
}

// Unwrap returns the matched value.
// Panics if no match was made.
func (m *Matcher) Unwrap() interface{} {
	if !m.matched {
		panic("pattern: attempted to unwrap unmatched value")
	}
	return m.value
}

// UnwrapOr returns the matched value or a default.
func (m *Matcher) UnwrapOr(defaultValue interface{}) interface{} {
	if m.matched {
		return m.value
	}
	return defaultValue
}

// UnwrapOrElse returns the matched value or computes a default.
func (m *Matcher) UnwrapOrElse(f func() interface{}) interface{} {
	if m.matched {
		return m.value
	}
	return f()
}

// MatchString creates a new StringMatcher for string pattern matching.
func MatchString(value string) *StringMatcher {
	return &StringMatcher{
		Matcher: Matcher{
			value:   value,
			matched: false,
		},
	}
}

// StringMatcher provides enhanced string pattern matching.
type StringMatcher struct {
	Matcher
}

// Prefix matches strings with a specific prefix.
func (m *StringMatcher) Prefix(prefix string, f func(string)) *StringMatcher {
	if m.matched {
		return m
	}

	if str, ok := m.value.(string); ok {
		if len(str) >= len(prefix) && str[:len(prefix)] == prefix {
			f(str)
			m.matched = true
		}
	}
	return m
}

// Suffix matches strings with a specific suffix.
func (m *StringMatcher) Suffix(suffix string, f func(string)) *StringMatcher {
	if m.matched {
		return m
	}

	if str, ok := m.value.(string); ok {
		if len(str) >= len(suffix) && str[len(str)-len(suffix):] == suffix {
			f(str)
			m.matched = true
		}
	}
	return m
}

// Contains matches strings containing a substring.
func (m *StringMatcher) Contains(substr string, f func(string)) *StringMatcher {
	if m.matched {
		return m
	}

	if str, ok := m.value.(string); ok {
		for i := 0; i <= len(str)-len(substr); i++ {
			if str[i:i+len(substr)] == substr {
				f(str)
				m.matched = true
				break
			}
		}
	}
	return m
}
