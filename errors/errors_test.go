// Package errors_test provides tests for the enhanced error handling utilities.
package errors_test

import (
	"fmt"
	"testing"

	"github.com/dongrv/rust-go/errors"
)

func TestNew(t *testing.T) {
	err := errors.New("test error")
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("error with value: %d", 42)
	expected := "error with value: 42"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestWrap(t *testing.T) {
	original := fmt.Errorf("original error")
	wrapped := errors.Wrap(original, "context")

	if wrapped.Error() != "context: original error" {
		t.Errorf("Expected 'context: original error', got '%s'", wrapped.Error())
	}

	// Test wrapping nil
	if errors.Wrap(nil, "context") != nil {
		t.Error("Wrap should return nil when wrapping nil")
	}
}

func TestWrapf(t *testing.T) {
	original := fmt.Errorf("original error")
	wrapped := errors.Wrapf(original, "context with %s", "value")

	expected := "context with value: original error"
	if wrapped.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, wrapped.Error())
	}
}

func TestWithContext(t *testing.T) {
	err := errors.New("test error").
		WithContext("key1", "value1").
		WithContext("key2", 42)

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got %v", err.Context["key1"])
	}
	if err.Context["key2"] != 42 {
		t.Errorf("Expected key2=42, got %v", err.Context["key2"])
	}
}

func TestWithContextMap(t *testing.T) {
	context := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	err := errors.New("test error").WithContextMap(context)

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got %v", err.Context["key1"])
	}
	if err.Context["key2"] != 42 {
		t.Errorf("Expected key2=42, got %v", err.Context["key2"])
	}
}

func TestResultOk(t *testing.T) {
	result := errors.Ok(42)
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.IsErr() {
		t.Error("Result should not be Err")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}
}

func TestResultErr(t *testing.T) {
	err := fmt.Errorf("test error")
	result := errors.Err[int](err)
	if result.IsOk() {
		t.Error("Result should not be Ok")
	}
	if !result.IsErr() {
		t.Error("Result should be Err")
	}
	if result.Error() != err {
		t.Errorf("Expected error '%v', got '%v'", err, result.Error())
	}
}

func TestResultMap(t *testing.T) {
	// Test map on Ok
	result := errors.Ok(21).Map(func(x int) int { return x * 2 })
	if !result.IsOk() {
		t.Error("Result should be Ok after Map")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}

	// Test map on Err
	err := fmt.Errorf("test error")
	result2 := errors.Err[int](err).Map(func(x int) int { return x * 2 })
	if !result2.IsErr() {
		t.Error("Result should still be Err after Map")
	}
}

func TestResultMapErr(t *testing.T) {
	// Test MapErr on Err
	err := fmt.Errorf("original error")
	result := errors.Err[int](err).MapErr(func(e error) error {
		return fmt.Errorf("wrapped: %v", e)
	})

	if !result.IsErr() {
		t.Error("Result should be Err")
	}
	if result.Error().Error() != "wrapped: original error" {
		t.Errorf("Expected 'wrapped: original error', got '%v'", result.Error())
	}

	// Test MapErr on Ok
	result2 := errors.Ok(42).MapErr(func(e error) error {
		return fmt.Errorf("should not be called")
	})
	if !result2.IsOk() {
		t.Error("Result should still be Ok after MapErr")
	}
}

func TestResultAndThen(t *testing.T) {
	result := errors.Ok(2).AndThen(func(x int) errors.Result[int] {
		return errors.Ok(x * 3)
	})

	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.Unwrap() != 6 {
		t.Errorf("Expected 6, got %v", result.Unwrap())
	}

	// Test AndThen on Err
	err := fmt.Errorf("test error")
	result2 := errors.Err[int](err).AndThen(func(x int) errors.Result[int] {
		return errors.Ok(x * 3)
	})
	if !result2.IsErr() {
		t.Error("Result should still be Err after AndThen")
	}
}

func TestResultOrElse(t *testing.T) {
	// Test OrElse on Err
	err := fmt.Errorf("test error")
	result := errors.Err[int](err).OrElse(func(e error) errors.Result[int] {
		return errors.Ok(42)
	})

	if !result.IsOk() {
		t.Error("Result should be Ok after OrElse")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}

	// Test OrElse on Ok
	result2 := errors.Ok(21).OrElse(func(e error) errors.Result[int] {
		return errors.Ok(999) // Should not be called
	})
	if !result2.IsOk() {
		t.Error("Result should still be Ok")
	}
	if result2.Unwrap() != 21 {
		t.Errorf("Expected 21, got %v", result2.Unwrap())
	}
}

func TestResultUnwrapOr(t *testing.T) {
	// Test UnwrapOr on Ok
	result := errors.Ok(42)
	if result.UnwrapOr(0) != 42 {
		t.Errorf("Expected 42, got %v", result.UnwrapOr(0))
	}

	// Test UnwrapOr on Err
	result2 := errors.Err[int](fmt.Errorf("error"))
	if result2.UnwrapOr(99) != 99 {
		t.Errorf("Expected 99, got %v", result2.UnwrapOr(99))
	}
}

func TestResultUnwrapOrElse(t *testing.T) {
	// Test UnwrapOrElse on Ok
	result := errors.Ok(42)
	value := result.UnwrapOrElse(func(e error) int {
		return 999 // Should not be called
	})
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	// Test UnwrapOrElse on Err
	result2 := errors.Err[int](fmt.Errorf("error"))
	value2 := result2.UnwrapOrElse(func(e error) int {
		return 99
	})
	if value2 != 99 {
		t.Errorf("Expected 99, got %v", value2)
	}
}

func TestResultExpect(t *testing.T) {
	// Test Expect on Ok
	result := errors.Ok(42)
	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not panic on Ok")
		}
	}()
	if result.Expect("should not panic") != 42 {
		t.Error("Expected 42")
	}

	// Test Expect on Err
	result2 := errors.Err[int](fmt.Errorf("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic on Err")
		}
	}()
	result2.Expect("this should panic")
}

func TestTry(t *testing.T) {
	// Test Try with no error
	result := errors.Try(42, nil)
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}

	// Test Try with error
	err := fmt.Errorf("test error")
	result2 := errors.Try(0, err)
	if !result2.IsErr() {
		t.Error("Result should be Err")
	}
	if result2.Error() != err {
		t.Errorf("Expected error '%v', got '%v'", err, result2.Error())
	}
}

func TestTryFunc(t *testing.T) {
	// Test TryFunc with no error
	result := errors.TryFunc(func() (int, error) {
		return 42, nil
	})
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}

	// Test TryFunc with error
	err := fmt.Errorf("test error")
	result2 := errors.TryFunc(func() (int, error) {
		return 0, err
	})
	if !result2.IsErr() {
		t.Error("Result should be Err")
	}
	if result2.Error() != err {
		t.Errorf("Expected error '%v', got '%v'", err, result2.Error())
	}
}

func TestRecover(t *testing.T) {
	// Test Recover without panic
	result := errors.Recover(func() int {
		return 42
	})
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.Unwrap() != 42 {
		t.Errorf("Expected 42, got %v", result.Unwrap())
	}

	// Test Recover with panic
	result2 := errors.Recover(func() int {
		panic("test panic")
	})
	if !result2.IsErr() {
		t.Error("Result should be Err after panic")
	}
	if result2.Error().Error() != "panic recovered: test panic" {
		t.Errorf("Expected 'panic recovered: test panic', got '%v'", result2.Error())
	}
}

func TestCombine(t *testing.T) {
	// Test Combine with all Ok
	results := []errors.Result[int]{
		errors.Ok(1),
		errors.Ok(2),
		errors.Ok(3),
	}
	combined := errors.Combine(results...)
	if !combined.IsOk() {
		t.Error("Combined result should be Ok")
	}
	values := combined.Unwrap()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", values)
	}

	// Test Combine with error
	results2 := []errors.Result[int]{
		errors.Ok(1),
		errors.Err[int](fmt.Errorf("error")),
		errors.Ok(3),
	}
	combined2 := errors.Combine(results2...)
	if !combined2.IsErr() {
		t.Error("Combined result should be Err")
	}
}

func TestFirstError(t *testing.T) {
	// Test with no errors
	results := []errors.Result[int]{
		errors.Ok(1),
		errors.Ok(2),
	}
	if errors.FirstError(results...) != nil {
		t.Error("Should return nil when no errors")
	}

	// Test with error
	err := fmt.Errorf("first error")
	results2 := []errors.Result[int]{
		errors.Ok(1),
		errors.Err[int](err),
		errors.Err[int](fmt.Errorf("second error")),
	}
	if errors.FirstError(results2...) != err {
		t.Error("Should return first error")
	}
}

func TestAllOk(t *testing.T) {
	// Test all Ok
	results := []errors.Result[int]{
		errors.Ok(1),
		errors.Ok(2),
	}
	if !errors.AllOk(results...) {
		t.Error("AllOk should return true when all are Ok")
	}

	// Test with error
	results2 := []errors.Result[int]{
		errors.Ok(1),
		errors.Err[int](fmt.Errorf("error")),
	}
	if errors.AllOk(results2...) {
		t.Error("AllOk should return false when any is Err")
	}
}

func TestAnyErr(t *testing.T) {
	// Test with no errors
	results := []errors.Result[int]{
		errors.Ok(1),
		errors.Ok(2),
	}
	if errors.AnyErr(results...) {
		t.Error("AnyErr should return false when no errors")
	}

	// Test with error
	results2 := []errors.Result[int]{
		errors.Ok(1),
		errors.Err[int](fmt.Errorf("error")),
	}
	if !errors.AnyErr(results2...) {
		t.Error("AnyErr should return true when any is Err")
	}
}

func TestErrorHandler(t *testing.T) {
	// Test Handle with no error
	handler := errors.Handle(nil).Then(func() error {
		return nil
	})
	if handler.Unwrap() != nil {
		t.Error("Should have no error")
	}

	// Test Handle with error
	err := fmt.Errorf("test error")
	handler2 := errors.Handle(err).Then(func() error {
		return nil // Should not be called
	})
	if handler2.Unwrap() != err {
		t.Errorf("Expected error '%v', got '%v'", err, handler2.Unwrap())
	}

	// Test If condition
	handler3 := errors.Handle(nil).If(true).Then(func() error {
		return fmt.Errorf("conditional error")
	})
	if handler3.Unwrap() == nil {
		t.Error("Should have error when condition is true")
	}

	handler4 := errors.Handle(nil).If(false).Then(func() error {
		return fmt.Errorf("should not be called")
	})
	// When condition is false, Then should not be called
	if handler4.Unwrap() != nil {
		t.Error("Should have no error when condition is false")
	}
}

func TestErrorHandlerMap(t *testing.T) {
	err := fmt.Errorf("original error")
	handler := errors.Handle(err).Map(func(e error) error {
		return fmt.Errorf("mapped: %v", e)
	})

	if handler.Unwrap().Error() != "mapped: original error" {
		t.Errorf("Expected 'mapped: original error', got '%v'", handler.Unwrap())
	}
}

func TestErrorHandlerOrElse(t *testing.T) {
	// Test OrElse with error
	err := fmt.Errorf("test error")
	newErr := errors.Handle(err).OrElse(func() error {
		return nil // Should not be called
	})
	if newErr != err {
		t.Errorf("Expected error '%v', got '%v'", err, newErr)
	}

	// Test OrElse without error
	newErr2 := errors.Handle(nil).OrElse(func() error {
		return fmt.Errorf("new error")
	})
	if newErr2.Error() != "new error" {
		t.Errorf("Expected 'new error', got '%v'", newErr2)
	}
}

func TestErrorHandlerMust(t *testing.T) {
	// Test Must without error
	handler := errors.Handle(nil)
	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not panic without error")
		}
	}()
	handler.Must()

	// Test Must with error
	handler2 := errors.Handle(fmt.Errorf("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic with error")
		}
	}()
	handler2.Must()
}

func TestErrorChain(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	err3 := fmt.Errorf("error 3")

	chain := errors.NewChain(err1, err2, err3)

	if chain.First() != err1 {
		t.Error("First should return first error")
	}
	if chain.Last() != err3 {
		t.Error("Last should return last error")
	}
	if len(chain) != 3 {
		t.Errorf("Expected chain length 3, got %d", len(chain))
	}

	// Test Append
	chain2 := chain.Append(fmt.Errorf("error 4"))
	if len(chain2) != 4 {
		t.Errorf("Expected chain length 4 after append, got %d", len(chain2))
	}
}

func TestFlatten(t *testing.T) {
	err1 := fmt.Errorf("error 1")

	// Create nested errors
	nested := errors.Wrap(errors.Wrap(err1, "wrap2"), "wrap1")

	chain := errors.Flatten(nested)
	if len(chain) != 3 {
		t.Errorf("Expected 3 errors in chain (wrap1, wrap2, and original), got %d", len(chain))
	}
}

func TestErrorString(t *testing.T) {
	err := errors.New("test error").
		WithContext("key", "value").
		WithContext("number", 42)

	str := err.String()
	// String() doesn't include "Error: " prefix anymore
	if !contains(str, "test error") {
		t.Error("String should contain error message")
	}
	if !contains(str, "key: value") {
		t.Error("String should contain context")
	}
	if !contains(str, "number: 42") {
		t.Error("String should contain context")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
