// Package immutable_test provides tests for the immutable data structures.
package immutable_test

import (
	"testing"

	"github.com/dongrv/rust-go/immutable"
)

func TestList(t *testing.T) {
	// Test EmptyList
	list := immutable.EmptyList[int]()
	if !list.IsEmpty() {
		t.Error("EmptyList should be empty")
	}
	if list.Size() != 0 {
		t.Errorf("EmptyList size should be 0, got %d", list.Size())
	}

	// Test ListOf
	list = immutable.ListOf(1, 2, 3, 4, 5)
	if list.IsEmpty() {
		t.Error("ListOf should not be empty")
	}
	if list.Size() != 5 {
		t.Errorf("Expected size 5, got %d", list.Size())
	}

	// Test Head
	if list.Head() != 1 {
		t.Errorf("Expected head 1, got %d", list.Head())
	}

	// Test Tail
	tail := list.Tail()
	if tail.Size() != 4 {
		t.Errorf("Expected tail size 4, got %d", tail.Size())
	}
	if tail.Head() != 2 {
		t.Errorf("Expected tail head 2, got %d", tail.Head())
	}

	// Test Cons
	newList := list.Cons(0)
	if newList.Size() != 6 {
		t.Errorf("Expected size 6 after Cons, got %d", newList.Size())
	}
	if newList.Head() != 0 {
		t.Errorf("Expected new head 0, got %d", newList.Head())
	}

	// Test Append
	list1 := immutable.ListOf(1, 2, 3)
	list2 := immutable.ListOf(4, 5, 6)
	appended := list1.Append(list2)
	if appended.Size() != 6 {
		t.Errorf("Expected appended size 6, got %d", appended.Size())
	}

	// Test Map
	mapped := list.Map(func(x int) int { return x * 2 })
	expected := []int{2, 4, 6, 8, 10}
	for i, v := range mapped.ToSlice() {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	// Test Filter
	filtered := list.Filter(func(x int) bool { return x%2 == 0 })
	if filtered.Size() != 2 {
		t.Errorf("Expected 2 even numbers, got %d", filtered.Size())
	}

	// Test FoldLeft
	sum := list.FoldLeft(0, func(acc interface{}, x int) interface{} {
		return acc.(int) + x
	}).(int)
	if sum != 15 {
		t.Errorf("Expected sum 15, got %d", sum)
	}

	// Test Reverse
	reversed := list.Reverse()
	reversedSlice := reversed.ToSlice()
	expectedReversed := []int{5, 4, 3, 2, 1}
	for i, v := range reversedSlice {
		if v != expectedReversed[i] {
			t.Errorf("Expected %d at index %d in reversed, got %d", expectedReversed[i], i, v)
		}
	}

	// Test ToSlice
	slice := list.ToSlice()
	if len(slice) != 5 {
		t.Errorf("Expected slice length 5, got %d", len(slice))
	}
	for i, v := range slice {
		if v != i+1 {
			t.Errorf("Expected %d at index %d, got %d", i+1, i, v)
		}
	}
}

func TestVector(t *testing.T) {
	// Test EmptyVector
	vector := immutable.EmptyVector[int]()
	if !vector.IsEmpty() {
		t.Error("EmptyVector should be empty")
	}
	if vector.Length() != 0 {
		t.Errorf("EmptyVector length should be 0, got %d", vector.Length())
	}

	// Test VectorOf
	vector = immutable.VectorOf(1, 2, 3, 4, 5)
	if vector.IsEmpty() {
		t.Error("VectorOf should not be empty")
	}
	if vector.Length() != 5 {
		t.Errorf("Expected length 5, got %d", vector.Length())
	}

	// Test Get
	for i := 0; i < 5; i++ {
		if vector.Get(i) != i+1 {
			t.Errorf("Expected %d at index %d, got %d", i+1, i, vector.Get(i))
		}
	}

	// Test Append
	appended := vector.Append(6)
	if appended.Length() != 6 {
		t.Errorf("Expected length 6 after Append, got %d", appended.Length())
	}
	if appended.Get(5) != 6 {
		t.Errorf("Expected 6 at index 5, got %d", appended.Get(5))
	}

	// Test Set
	updated := vector.Set(2, 99)
	if updated.Get(2) != 99 {
		t.Errorf("Expected 99 at index 2 after Set, got %d", updated.Get(2))
	}
	// Original should be unchanged
	if vector.Get(2) != 3 {
		t.Errorf("Original vector should be unchanged, expected 3 at index 2, got %d", vector.Get(2))
	}

	// Test Map
	mapped := vector.Map(func(x int) int { return x * 2 })
	expected := []int{2, 4, 6, 8, 10}
	for i := 0; i < mapped.Length(); i++ {
		if mapped.Get(i) != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, mapped.Get(i))
		}
	}

	// Test Filter
	filtered := vector.Filter(func(x int) bool { return x%2 == 0 })
	if filtered.Length() != 2 {
		t.Errorf("Expected 2 even numbers, got %d", filtered.Length())
	}

	// Test ToSlice
	slice := vector.ToSlice()
	if len(slice) != 5 {
		t.Errorf("Expected slice length 5, got %d", len(slice))
	}
	for i, v := range slice {
		if v != i+1 {
			t.Errorf("Expected %d at index %d, got %d", i+1, i, v)
		}
	}
}

func TestMap(t *testing.T) {
	// Test EmptyMap
	m := immutable.EmptyMap[string, int]()
	if !m.IsEmpty() {
		t.Error("EmptyMap should be empty")
	}
	if m.Size() != 0 {
		t.Errorf("EmptyMap size should be 0, got %d", m.Size())
	}

	// Test MapOf
	m = immutable.MapOf(
		immutable.PairOf("one", 1),
		immutable.PairOf("two", 2),
		immutable.PairOf("three", 3),
	)
	if m.IsEmpty() {
		t.Error("MapOf should not be empty")
	}
	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}

	// Test Get
	if val, ok := m.Get("one"); !ok || val != 1 {
		t.Errorf("Expected (1, true) for key 'one', got (%d, %v)", val, ok)
	}
	if val, ok := m.Get("two"); !ok || val != 2 {
		t.Errorf("Expected (2, true) for key 'two', got (%d, %v)", val, ok)
	}
	if _, ok := m.Get("four"); ok {
		t.Error("Expected false for non-existent key 'four'")
	}

	// Test Set
	m2 := m.Set("four", 4)
	if m2.Size() != 4 {
		t.Errorf("Expected size 4 after Set, got %d", m2.Size())
	}
	if val, ok := m2.Get("four"); !ok || val != 4 {
		t.Errorf("Expected (4, true) for key 'four', got (%d, %v)", val, ok)
	}
	// Original should be unchanged
	if m.Size() != 3 {
		t.Errorf("Original map should be unchanged, expected size 3, got %d", m.Size())
	}

	// Test Update existing key
	m3 := m.Set("one", 100)
	if val, ok := m3.Get("one"); !ok || val != 100 {
		t.Errorf("Expected (100, true) for updated key 'one', got (%d, %v)", val, ok)
	}
	if m3.Size() != 3 {
		t.Errorf("Size should remain 3 after update, got %d", m3.Size())
	}

	// Test Delete
	m4 := m.Delete("two")
	if m4.Size() != 2 {
		t.Errorf("Expected size 2 after Delete, got %d", m4.Size())
	}
	if _, ok := m4.Get("two"); ok {
		t.Error("Key 'two' should be deleted")
	}
	// Original should be unchanged
	if _, ok := m.Get("two"); !ok {
		t.Error("Original map should still have key 'two'")
	}

	// Test Contains
	if !m.Contains("one") {
		t.Error("Map should contain key 'one'")
	}
	if m.Contains("four") {
		t.Error("Map should not contain key 'four'")
	}

	// Test ForEach
	sum := 0
	count := 0
	m.ForEach(func(key string, value int) {
		sum += value
		count++
	})
	if count != 3 {
		t.Errorf("Expected 3 items from ForEach, got %d", count)
	}
	if sum != 6 {
		t.Errorf("Expected sum 6 from ForEach, got %d", sum)
	}

	// Test Map (transform values)
	mapped := m.Map(func(value int) int { return value * 10 })
	// Note: Map creates new map, so we need to check all values
	expectedMap := map[string]int{
		"one":   10,
		"two":   20,
		"three": 30,
	}
	for key, expected := range expectedMap {
		if val, ok := mapped.Get(key); !ok || val != expected {
			t.Errorf("Expected (%d, true) for mapped key '%s', got (%d, %v)", expected, key, val, ok)
		}
	}

	// Test Filter
	filtered := m.Filter(func(key string, value int) bool { return value > 1 })
	if filtered.Size() != 2 {
		t.Errorf("Expected 2 elements after filter, got %d", filtered.Size())
	}
	// Check that filtered map contains only values > 1
	filtered.ForEach(func(key string, value int) {
		if value <= 1 {
			t.Errorf("Filtered map should not contain value %d", value)
		}
	})

	// Test Keys
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Test Values
	values := m.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test ToSlice
	pairs := m.ToSlice()
	if len(pairs) != 3 {
		t.Errorf("Expected 3 pairs, got %d", len(pairs))
	}
}

func TestSet(t *testing.T) {
	// Test EmptySet
	s := immutable.EmptySet[int]()
	if !s.IsEmpty() {
		t.Error("EmptySet should be empty")
	}
	if s.Size() != 0 {
		t.Errorf("EmptySet size should be 0, got %d", s.Size())
	}

	// Test SetOf
	s = immutable.SetOf(1, 2, 3, 4, 5)
	if s.IsEmpty() {
		t.Error("SetOf should not be empty")
	}
	if s.Size() != 5 {
		t.Errorf("Expected size 5, got %d", s.Size())
	}

	// Test Contains
	if !s.Contains(3) {
		t.Error("Set should contain 3")
	}
	if s.Contains(6) {
		t.Error("Set should not contain 6")
	}

	// Test Add
	s2 := s.Add(6)
	if s2.Size() != 6 {
		t.Errorf("Expected size 6 after Add, got %d", s2.Size())
	}
	if !s2.Contains(6) {
		t.Error("Added set should contain 6")
	}
	// Original should be unchanged
	if s.Size() != 5 {
		t.Errorf("Original set size should be 5, got %d", s.Size())
	}

	// Test Remove
	s3 := s.Remove(3)
	if s3.Size() != 4 {
		t.Errorf("Expected size 4 after Remove, got %d", s3.Size())
	}
	if s3.Contains(3) {
		t.Error("Removed set should not contain 3")
	}
	// Original should be unchanged
	if s.Size() != 5 {
		t.Errorf("Original set size should be 5, got %d", s.Size())
	}

	// Test Union
	s1 := immutable.SetOf(1, 2, 3)
	s2 = immutable.SetOf(3, 4, 5)
	union := s1.Union(s2)
	if union.Size() != 5 {
		t.Errorf("Expected union size 5, got %d", union.Size())
	}
	for i := 1; i <= 5; i++ {
		if !union.Contains(i) {
			t.Errorf("Union should contain %d", i)
		}
	}

	// Test Intersection
	intersection := s1.Intersection(s2)
	if intersection.Size() != 1 {
		t.Errorf("Expected intersection size 1, got %d", intersection.Size())
	}
	if !intersection.Contains(3) {
		t.Error("Intersection should contain 3")
	}

	// Test Difference
	difference := s1.Difference(s2)
	if difference.Size() != 2 {
		t.Errorf("Expected difference size 2, got %d", difference.Size())
	}
	if !difference.Contains(1) || !difference.Contains(2) {
		t.Error("Difference should contain 1 and 2")
	}
	if difference.Contains(3) {
		t.Error("Difference should not contain 3")
	}

	// Test ForEach
	sum := 0
	count := 0
	s.ForEach(func(value int) {
		sum += value
		count++
	})
	if count != 5 {
		t.Errorf("Expected 5 items from ForEach, got %d", count)
	}
	if sum != 15 {
		t.Errorf("Expected sum 15 from ForEach, got %d", sum)
	}

	// Test ToSlice
	slice := s.ToSlice()
	if len(slice) != 5 {
		t.Errorf("Expected slice length 5, got %d", len(slice))
	}
}

func TestImmutableProperty(t *testing.T) {
	// Test that operations return new instances
	list1 := immutable.ListOf(1, 2, 3)
	list2 := list1.Cons(0)

	if list1.Size() != 3 {
		t.Error("Original list should be unchanged")
	}
	if list2.Size() != 4 {
		t.Error("New list should have new element")
	}

	// Test vector immutability
	vector1 := immutable.VectorOf(1, 2, 3)
	vector2 := vector1.Set(1, 99)

	if vector1.Get(1) != 2 {
		t.Error("Original vector should be unchanged")
	}
	if vector2.Get(1) != 99 {
		t.Error("New vector should have updated value")
	}

	// Test map immutability
	map1 := immutable.MapOf(immutable.PairOf("a", 1))
	map2 := map1.Set("b", 2)

	if map1.Size() != 1 {
		t.Error("Original map should be unchanged")
	}
	if map2.Size() != 2 {
		t.Error("New map should have new entry")
	}

	// Test set immutability
	set1 := immutable.SetOf(1, 2, 3)
	set2 := set1.Add(4)

	if set1.Size() != 3 {
		t.Error("Original set should be unchanged")
	}
	if set2.Size() != 4 {
		t.Error("New set should have new element")
	}
}

func TestStringRepresentation(t *testing.T) {
	// Test List string
	list := immutable.ListOf(1, 2, 3)
	str := list.String()
	if str != "List[1, 2, 3]" {
		t.Errorf("Expected 'List[1, 2, 3]', got '%s'", str)
	}

	// Test Vector string
	vector := immutable.VectorOf(1, 2, 3)
	str = vector.String()
	if str != "Vector[1, 2, 3]" {
		t.Errorf("Expected 'Vector[1, 2, 3]', got '%s'", str)
	}

	// Test Map string
	m := immutable.MapOf(immutable.PairOf("a", 1), immutable.PairOf("b", 2))
	str = m.String()
	// Order may vary, just check it contains the pairs
	if !contains(str, "a: 1") || !contains(str, "b: 2") {
		t.Errorf("String should contain both pairs, got '%s'", str)
	}

	// Test Set string
	s := immutable.SetOf(1, 2, 3)
	str = s.String()
	// Order may vary, just check it contains the numbers
	if !contains(str, "1") || !contains(str, "2") || !contains(str, "3") {
		t.Errorf("String should contain all numbers, got '%s'", str)
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

func TestEdgeCases(t *testing.T) {
	// Test empty list operations
	emptyList := immutable.EmptyList[int]()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Head() on empty list should panic")
		}
	}()
	emptyList.Head()

	// Test vector bounds
	vector := immutable.VectorOf(1, 2, 3)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Get() with out of bounds index should panic")
		}
	}()
	vector.Get(5)

	// Test map with nil values
	m := immutable.EmptyMap[string, *int]()
	val := 42
	m2 := m.Set("key", &val)
	if retrieved, ok := m2.Get("key"); !ok || *retrieved != 42 {
		t.Error("Map should handle pointer values")
	}
}

func BenchmarkListOperations(b *testing.B) {
	list := immutable.EmptyList[int]()
	for i := 0; i < b.N; i++ {
		list = list.Cons(i)
	}
}
