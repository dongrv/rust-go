// Package pattern_test provides tests for the pattern matching library
package pattern_test

import (
	"fmt"
	"testing"

	"github.com/dongrv/rust-go"
	"github.com/dongrv/rust-go/pattern"
)

// TestMatchOptionSome tests matching Some values
func TestMatchOptionSome(t *testing.T) {
	t.Run("Some with value", func(t *testing.T) {
		value := rust.Some(42)
		called := false

		pattern.Match(value).
			Some(func(x int) {
				if x != 42 {
					t.Errorf("Expected 42, got %d", x)
				}
				called = true
			}).
			None(func() {
				t.Error("Should not call None for Some(42)")
			})

		if !called {
			t.Error("Some handler was not called")
		}
	})

	t.Run("Some with chaining", func(t *testing.T) {
		value := rust.Some("hello")
		someCalled := false
		noneCalled := false

		pattern.Match(value).
			Some(func(s string) {
				if s != "hello" {
					t.Errorf("Expected 'hello', got %s", s)
				}
				someCalled = true
			}).
			None(func() {
				noneCalled = true
			})

		if !someCalled {
			t.Error("Some handler was not called")
		}
		if noneCalled {
			t.Error("None handler should not be called")
		}
	})
}

// TestMatchOptionNone tests matching None values
func TestMatchOptionNone(t *testing.T) {
	t.Run("None value", func(t *testing.T) {
		value := rust.None[int]()
		called := false

		pattern.Match(value).
			Some(func(x int) {
				t.Error("Should not call Some for None")
			}).
			None(func() {
				called = true
			})

		if !called {
			t.Error("None handler was not called")
		}
	})
}

// TestMatchResultOk tests matching Ok values
func TestMatchResultOk(t *testing.T) {
	t.Run("Ok with value", func(t *testing.T) {
		value := rust.Ok[int, string](42)
		okCalled := false
		errCalled := false

		pattern.Match(value).
			Ok(func(x int) {
				if x != 42 {
					t.Errorf("Expected 42, got %d", x)
				}
				okCalled = true
			}).
			Err(func(err string) {
				errCalled = true
			})

		if !okCalled {
			t.Error("Ok handler was not called")
		}
		if errCalled {
			t.Error("Err handler should not be called")
		}
	})
}

// TestMatchResultErr tests matching Err values
func TestMatchResultErr(t *testing.T) {
	t.Run("Err with error", func(t *testing.T) {
		value := rust.Err[int, string]("something went wrong")
		okCalled := false
		errCalled := false

		pattern.Match(value).
			Ok(func(x int) {
				okCalled = true
			}).
			Err(func(err string) {
				if err != "something went wrong" {
					t.Errorf("Expected 'something went wrong', got %s", err)
				}
				errCalled = true
			})

		if okCalled {
			t.Error("Ok handler should not be called")
		}
		if !errCalled {
			t.Error("Err handler was not called")
		}
	})
}

// TestMatchValue tests matching specific values
func TestMatchValue(t *testing.T) {
	t.Run("Exact value match", func(t *testing.T) {
		value := 42
		exactCalled := false
		defaultCalled := false

		pattern.Match(value).
			Value(42, func() {
				exactCalled = true
			}).
			Value(100, func() {
				t.Error("Should not match 100")
			}).
			Default(func() {
				defaultCalled = true
			})

		if !exactCalled {
			t.Error("Exact value handler was not called")
		}
		if defaultCalled {
			t.Error("Default handler should not be called")
		}
	})

	t.Run("Multiple value matches", func(t *testing.T) {
		value := "hello"
		helloCalled := false
		worldCalled := false

		pattern.Match(value).
			Value("hello", func() {
				helloCalled = true
			}).
			Value("world", func() {
				worldCalled = true
			})

		if !helloCalled {
			t.Error("'hello' handler was not called")
		}
		if worldCalled {
			t.Error("'world' handler should not be called")
		}
	})
}

// TestMatchType tests matching by type
func TestMatchType(t *testing.T) {
	t.Run("Type match with interface", func(t *testing.T) {
		var value interface{} = "hello"
		stringCalled := false
		intCalled := false

		pattern.Match(value).
			Type(func(s string) {
				if s != "hello" {
					t.Errorf("Expected 'hello', got %s", s)
				}
				stringCalled = true
			}).
			Type(func(i int) {
				intCalled = true
			}).
			Default(func() {
				t.Error("Default should not be called")
			})

		if !stringCalled {
			t.Error("String type handler was not called")
		}
		if intCalled {
			t.Error("Int type handler should not be called")
		}
	})

	t.Run("Type match with int", func(t *testing.T) {
		var value interface{} = 42
		intCalled := false

		pattern.Match(value).
			Type(func(i int) {
				if i != 42 {
					t.Errorf("Expected 42, got %d", i)
				}
				intCalled = true
			})

		if !intCalled {
			t.Error("Int type handler was not called")
		}
	})
}

// TestMatchPredicate tests matching with custom predicates
func TestMatchPredicate(t *testing.T) {
	t.Run("Predicate match", func(t *testing.T) {
		value := 42
		evenCalled := false
		oddCalled := false

		pattern.Match(value).
			Predicate(func(x int) bool { return x%2 == 0 }, func() {
				evenCalled = true
			}).
			Predicate(func(x int) bool { return x%2 != 0 }, func() {
				oddCalled = true
			})

		if !evenCalled {
			t.Error("Even predicate handler was not called")
		}
		if oddCalled {
			t.Error("Odd predicate handler should not be called")
		}
	})

	t.Run("Multiple predicates", func(t *testing.T) {
		value := "hello"
		shortCalled := false
		longCalled := false

		pattern.Match(value).
			Predicate(func(s string) bool { return len(s) < 10 }, func() {
				shortCalled = true
			}).
			Predicate(func(s string) bool { return len(s) >= 10 }, func() {
				longCalled = true
			})

		if !shortCalled {
			t.Error("Short predicate handler was not called")
		}
		if longCalled {
			t.Error("Long predicate handler should not be called")
		}
	})
}

// TestMatchDefault tests the default case
func TestMatchDefault(t *testing.T) {
	t.Run("Default case when no match", func(t *testing.T) {
		value := 100
		defaultCalled := false

		pattern.Match(value).
			Value(42, func() {
				t.Error("Should not match 42")
			}).
			Default(func() {
				defaultCalled = true
			})

		if !defaultCalled {
			t.Error("Default handler was not called")
		}
	})

	t.Run("Default not called when match found", func(t *testing.T) {
		value := 42
		valueCalled := false
		defaultCalled := false

		pattern.Match(value).
			Value(42, func() {
				valueCalled = true
			}).
			Default(func() {
				defaultCalled = true
			})

		if !valueCalled {
			t.Error("Value handler was not called")
		}
		if defaultCalled {
			t.Error("Default handler should not be called")
		}
	})
}

// TestMatchExhaustive tests exhaustive matching
func TestMatchExhaustive(t *testing.T) {
	t.Run("Exhaustive match success", func(t *testing.T) {
		value := rust.Some(42)

		defer func() {
			if r := recover(); r != nil {
				t.Error("Should not panic on exhaustive match")
			}
		}()

		pattern.Match(value).
			Some(func(x int) {
				// Do nothing
			}).
			None(func() {
				// Do nothing
			}).
			Exhaustive()
	})

	t.Run("Exhaustive match panic", func(t *testing.T) {
		value := 42

		defer func() {
			if r := recover(); r == nil {
				t.Error("Should panic on non-exhaustive match")
			}
		}()

		pattern.Match(value).
			Value(100, func() {
				// This won't match
			}).
			Exhaustive()
	})
}

// TestMatchMap tests mapping matched values
func TestMatchMap(t *testing.T) {
	t.Run("Map with transformation", func(t *testing.T) {
		value := 42
		result := pattern.Match(value).
			Value(42, func() {
				// Match first
			}).
			Map(func(x int) string {
				return fmt.Sprintf("Value: %d", x)
			}).
			UnwrapOr("default")

		expected := "Value: 42"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("Map with interface value", func(t *testing.T) {
		var value interface{} = 42
		result := pattern.Match(value).
			Type(func(x int) {
				// Match first
			}).
			Map(func(x int) string {
				return fmt.Sprintf("Number: %d", x)
			}).
			UnwrapOr("unknown")

		expected := "Number: 42"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

// TestMatchUnwrap tests unwrapping matched values
func TestMatchUnwrap(t *testing.T) {
	t.Run("Unwrap with match", func(t *testing.T) {
		value := 42
		result := pattern.Match(value).
			Value(42, func() {
				// Match found
			}).
			Unwrap()

		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("Unwrap panic without match", func(t *testing.T) {
		value := 42

		defer func() {
			if r := recover(); r == nil {
				t.Error("Should panic when unwrapping without match")
			}
		}()

		pattern.Match(value).
			Value(100, func() {
				// This won't match
			}).
			Unwrap()
	})

	t.Run("UnwrapOr with match", func(t *testing.T) {
		value := 42
		result := pattern.Match(value).
			Value(42, func() {
				// Match found
			}).
			UnwrapOr(100)

		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("UnwrapOr without match", func(t *testing.T) {
		value := 42
		result := pattern.Match(value).
			Value(100, func() {
				// This won't match
			}).
			UnwrapOr(999)

		if result != 999 {
			t.Errorf("Expected 999, got %v", result)
		}
	})

	t.Run("UnwrapOrElse without match", func(t *testing.T) {
		value := 42
		result := pattern.Match(value).
			Value(100, func() {
				// This won't match
			}).
			UnwrapOrElse(func() interface{} {
				return 777
			})

		if result != 777 {
			t.Errorf("Expected 777, got %v", result)
		}
	})
}

// TestMatchString tests string-specific pattern matching
func TestMatchString(t *testing.T) {
	t.Run("String prefix match", func(t *testing.T) {
		value := "hello world"
		prefixCalled := false

		pattern.MatchString(value).
			Prefix("hello", func(s string) {
				if s != "hello world" {
					t.Errorf("Expected 'hello world', got %s", s)
				}
				prefixCalled = true
			}).
			Default(func() {
				t.Error("Default should not be called")
			})

		if !prefixCalled {
			t.Error("Prefix handler was not called")
		}
	})

	t.Run("String suffix match", func(t *testing.T) {
		value := "hello world"
		suffixCalled := false

		pattern.MatchString(value).
			Suffix("world", func(s string) {
				if s != "hello world" {
					t.Errorf("Expected 'hello world', got %s", s)
				}
				suffixCalled = true
			})

		if !suffixCalled {
			t.Error("Suffix handler was not called")
		}
	})

	t.Run("String contains match", func(t *testing.T) {
		value := "hello world"
		containsCalled := false

		pattern.MatchString(value).
			Contains("lo wo", func(s string) {
				if s != "hello world" {
					t.Errorf("Expected 'hello world', got %s", s)
				}
				containsCalled = true
			})

		if !containsCalled {
			t.Error("Contains handler was not called")
		}
	})
}

// TestComplexPatterns tests complex pattern matching scenarios
func TestComplexPatterns(t *testing.T) {
	t.Run("Nested option matching", func(t *testing.T) {
		outer := rust.Some(rust.Some(42))
		called := false

		pattern.Match(outer).
			Some(func(inner rust.Option[int]) {
				pattern.Match(inner).
					Some(func(x int) {
						if x != 42 {
							t.Errorf("Expected 42, got %d", x)
						}
						called = true
					}).
					None(func() {
						t.Error("Inner should not be None")
					})
			}).
			None(func() {
				t.Error("Outer should not be None")
			})

		if !called {
			t.Error("Inner Some handler was not called")
		}
	})

	t.Run("Mixed pattern types", func(t *testing.T) {
		value := 42
		patternsMatched := 0

		pattern.Match(value).
			Value(42, func() {
				patternsMatched++
			}).
			Predicate(func(x int) bool { return x > 0 }, func() {
				patternsMatched++
			})

		// Only the first matching pattern should execute
		if patternsMatched != 1 {
			t.Errorf("Expected 1 pattern to match, got %d", patternsMatched)
		}
	})
}

// BenchmarkMatch tests performance of pattern matching
func BenchmarkMatch(b *testing.B) {
	value := rust.Some(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pattern.Match(value).
			Some(func(x int) {
				_ = x * 2
			}).
			None(func() {
				// Do nothing
			})
	}
}

// ExampleMatch demonstrates basic pattern matching
func ExampleMatch() {
	value := rust.Some(42)

	pattern.Match(value).
		Some(func(x int) {
			fmt.Printf("Got value: %d\n", x)
		}).
		None(func() {
			fmt.Println("Got nothing")
		})
	// Output:
	// Got value: 42
}
