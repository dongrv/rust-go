package rust_test

import (
	"fmt"
	"testing"

	. "github.com/dongrv/rust-go"
)

func TestOption(t *testing.T) {
	t.Run("Some creation", func(t *testing.T) {
		opt := Some(42)
		if !opt.IsSome() {
			t.Error("Expected Some to be IsSome")
		}
		if opt.IsNone() {
			t.Error("Expected Some not to be IsNone")
		}
		if opt.Unwrap() != 42 {
			t.Errorf("Expected Unwrap to return 42, got %v", opt.Unwrap())
		}
	})

	t.Run("None creation", func(t *testing.T) {
		opt := None[int]()
		if opt.IsSome() {
			t.Error("Expected None not to be IsSome")
		}
		if !opt.IsNone() {
			t.Error("Expected None to be IsNone")
		}
	})

	t.Run("UnwrapOr", func(t *testing.T) {
		some := Some(42)
		if some.UnwrapOr(100) != 42 {
			t.Error("Expected UnwrapOr to return Some value")
		}

		none := None[int]()
		if none.UnwrapOr(100) != 100 {
			t.Error("Expected UnwrapOr to return default for None")
		}
	})

	t.Run("MapOption", func(t *testing.T) {
		some := Some(21)
		mapped := MapOption(some, func(x int) int { return x * 2 })
		if mapped.UnwrapOr(0) != 42 {
			t.Error("Expected MapOption to transform value")
		}

		none := None[int]()
		mappedNone := MapOption(none, func(x int) int { return x * 2 })
		if mappedNone.IsSome() {
			t.Error("Expected MapOption on None to return None")
		}
	})

	t.Run("AndThenOption", func(t *testing.T) {
		some := Some(2)
		result := AndThenOption(some, func(x int) Option[int] {
			return Some(x * 21)
		})
		if result.UnwrapOr(0) != 42 {
			t.Error("Expected AndThenOption to chain operations")
		}

		none := None[int]()
		resultNone := AndThenOption(none, func(x int) Option[int] {
			return Some(x * 21)
		})
		if resultNone.IsSome() {
			t.Error("Expected AndThenOption on None to return None")
		}
	})
}

func TestResult(t *testing.T) {
	t.Run("Ok creation", func(t *testing.T) {
		res := Ok[int, string](42)
		if !res.IsOk() {
			t.Error("Expected Ok to be IsOk")
		}
		if res.IsErr() {
			t.Error("Expected Ok not to be IsErr")
		}
		if res.Unwrap() != 42 {
			t.Errorf("Expected Unwrap to return 42, got %v", res.Unwrap())
		}
	})

	t.Run("Err creation", func(t *testing.T) {
		res := Err[int, string]("error")
		if res.IsOk() {
			t.Error("Expected Err not to be IsOk")
		}
		if !res.IsErr() {
			t.Error("Expected Err to be IsErr")
		}
		if res.UnwrapErr() != "error" {
			t.Errorf("Expected UnwrapErr to return 'error', got %v", res.UnwrapErr())
		}
	})

	t.Run("MapResult", func(t *testing.T) {
		ok := Ok[int, string](21)
		mapped := MapResult(ok, func(x int) int { return x * 2 })
		if mapped.UnwrapOr(0) != 42 {
			t.Error("Expected MapResult to transform Ok value")
		}

		err := Err[int, string]("error")
		mappedErr := MapResult(err, func(x int) int { return x * 2 })
		if mappedErr.IsOk() {
			t.Error("Expected MapResult on Err to return Err")
		}
	})

	t.Run("AndThenResult", func(t *testing.T) {
		ok := Ok[int, string](2)
		result := AndThenResult(ok, func(x int) Result[int, string] {
			return Ok[int, string](x * 21)
		})
		if result.UnwrapOr(0) != 42 {
			t.Error("Expected AndThenResult to chain operations")
		}

		err := Err[int, string]("error")
		resultErr := AndThenResult(err, func(x int) Result[int, string] {
			return Ok[int, string](x * 21)
		})
		if resultErr.IsOk() {
			t.Error("Expected AndThenResult on Err to return Err")
		}
	})

	t.Run("UnwrapOrElse", func(t *testing.T) {
		ok := Ok[int, string](42)
		if ok.UnwrapOrElse(func(e string) int { return 100 }) != 42 {
			t.Error("Expected UnwrapOrElse to return Ok value")
		}

		err := Err[int, string]("error")
		if err.UnwrapOrElse(func(e string) int { return len(e) }) != 5 {
			t.Error("Expected UnwrapOrElse to compute from error")
		}
	})
}

func TestIterator(t *testing.T) {
	t.Run("Basic iteration", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		iter := Iter(slice)

		var sum int
		for {
			next := iter.Next()
			if next.IsNone() {
				break
			}
			sum += next.Unwrap()
		}

		if sum != 15 {
			t.Errorf("Expected sum 15, got %d", sum)
		}
	})

	t.Run("Map and Collect", func(t *testing.T) {
		slice := []int{1, 2, 3}
		result := Collect(Map(Iter(slice), func(x int) int { return x * 2 }))

		expected := []int{2, 4, 6}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Filter", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := Collect(Filter(Iter(slice), func(x int) bool { return x%2 == 0 }))

		expected := []int{2, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Fold", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		sum := Fold(Iter(slice), 0, func(acc, x int) int { return acc + x })

		if sum != 15 {
			t.Errorf("Expected sum 15, got %d", sum)
		}
	})

	t.Run("Take and Skip", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := Collect(Take(Skip(Iter(slice), 2), 2))

		expected := []int{3, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})
}

func TestChainable(t *testing.T) {
	t.Run("From and Collect", func(t *testing.T) {
		slice := []int{1, 2, 3}
		chain := From(slice)
		result := chain.Collect()

		if len(result) != len(slice) {
			t.Errorf("Expected length %d, got %d", len(slice), len(result))
		}
		for i, v := range slice {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Map and Filter", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := From(slice).
			Filter(func(x int) bool { return x%2 == 0 }).
			Map(func(x int) int { return x * 3 }).
			Collect()

		expected := []int{6, 12}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Reduce", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		sum := From(slice).
			Reduce(func(a, b int) int { return a + b }).
			UnwrapOr(0)

		if sum != 15 {
			t.Errorf("Expected sum 15, got %d", sum)
		}

		empty := []int{}
		sumEmpty := From(empty).
			Reduce(func(a, b int) int { return a + b })
		if sumEmpty.IsSome() {
			t.Error("Expected Reduce on empty slice to return None")
		}
	})

	t.Run("All and Any", func(t *testing.T) {
		slice := []int{2, 4, 6, 8}
		allEven := From(slice).All(func(x int) bool { return x%2 == 0 })
		if !allEven {
			t.Error("Expected All to return true for all even numbers")
		}

		anyOdd := From(slice).Any(func(x int) bool { return x%2 == 1 })
		if anyOdd {
			t.Error("Expected Any to return false for no odd numbers")
		}
	})

	t.Run("Take and Skip", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := From(slice).
			Skip(1).
			Take(3).
			Collect()

		expected := []int{2, 3, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := From(slice).Reverse().Collect()

		expected := []int{5, 4, 3, 2, 1}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Unique", func(t *testing.T) {
		slice := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
		result := From(slice).Unique().Collect()

		expected := []int{1, 2, 3, 4}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Partition", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		trueVals, falseVals := From(slice).Partition(func(x int) bool { return x > 3 })

		expectedTrue := []int{4, 5}
		expectedFalse := []int{1, 2, 3}

		trueResult := trueVals.Collect()
		falseResult := falseVals.Collect()

		if len(trueResult) != len(expectedTrue) {
			t.Errorf("Expected true length %d, got %d", len(expectedTrue), len(trueResult))
		}
		for i, v := range expectedTrue {
			if trueResult[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, trueResult[i])
			}
		}

		if len(falseResult) != len(expectedFalse) {
			t.Errorf("Expected false length %d, got %d", len(expectedFalse), len(falseResult))
		}
		for i, v := range expectedFalse {
			if falseResult[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, falseResult[i])
			}
		}
	})
}

func TestRange(t *testing.T) {
	t.Run("Range iterator", func(t *testing.T) {
		result := Collect(Range(1, 6, 1))
		expected := []int{1, 2, 3, 4, 5}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})

	t.Run("Range with step", func(t *testing.T) {
		result := Collect(Range(0, 10, 2))
		expected := []int{0, 2, 4, 6, 8}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
			}
		}
	})
}

func TestGenerate(t *testing.T) {
	t.Run("Generate sequence", func(t *testing.T) {
		result := Generate(5, func(i int) string {
			return fmt.Sprintf("item-%d", i+1)
		}).Collect()

		expected := []string{"item-1", "item-2", "item-3", "item-4", "item-5"}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
			}
		}
	})
}

func TestOnceAndRepeat(t *testing.T) {
	t.Run("Once iterator", func(t *testing.T) {
		result := Collect(Once("hello"))
		expected := []string{"hello"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
			}
		}
	})

	t.Run("Repeat with Take", func(t *testing.T) {
		result := Collect(Take(Repeat("loop"), 3))
		expected := []string{"loop", "loop", "loop"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
			}
		}
	})
}
