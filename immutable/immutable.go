// Package immutable provides persistent immutable data structures inspired by functional programming.
// These data structures are thread-safe, support efficient updates, and enable pure functional programming patterns.
package immutable

import (
	"fmt"
	"strings"
)

// List is a persistent immutable singly-linked list.
type List[T any] struct {
	head *listNode[T]
	size int
}

type listNode[T any] struct {
	value T
	next  *listNode[T]
}

// EmptyList creates an empty list.
func EmptyList[T any]() *List[T] {
	return &List[T]{head: nil, size: 0}
}

// ListOf creates a list from the given values.
func ListOf[T any](values ...T) *List[T] {
	if len(values) == 0 {
		return EmptyList[T]()
	}

	// Build list in reverse for efficiency
	var head *listNode[T]
	for i := len(values) - 1; i >= 0; i-- {
		head = &listNode[T]{
			value: values[i],
			next:  head,
		}
	}

	return &List[T]{
		head: head,
		size: len(values),
	}
}

// Cons adds an element to the front of the list.
// Returns a new list with the element added.
func (l *List[T]) Cons(value T) *List[T] {
	return &List[T]{
		head: &listNode[T]{
			value: value,
			next:  l.head,
		},
		size: l.size + 1,
	}
}

// Head returns the first element of the list.
// Panics if the list is empty.
func (l *List[T]) Head() T {
	if l.IsEmpty() {
		panic("List.Head: empty list")
	}
	return l.head.value
}

// Tail returns a new list without the first element.
// Returns empty list if the list is empty or has only one element.
func (l *List[T]) Tail() *List[T] {
	if l.IsEmpty() {
		return l
	}
	return &List[T]{
		head: l.head.next,
		size: l.size - 1,
	}
}

// IsEmpty returns true if the list is empty.
func (l *List[T]) IsEmpty() bool {
	return l.head == nil
}

// Size returns the number of elements in the list.
func (l *List[T]) Size() int {
	return l.size
}

// Append appends another list to this list.
// Returns a new list containing all elements.
func (l *List[T]) Append(other *List[T]) *List[T] {
	if l.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return l
	}

	// Reverse the first list to build the result
	reversed := EmptyList[T]()
	for node := l.head; node != nil; node = node.next {
		reversed = reversed.Cons(node.value)
	}

	// Append the second list
	result := other
	for node := reversed.head; node != nil; node = node.next {
		result = result.Cons(node.value)
	}

	return result
}

// Map applies a function to each element and returns a new list.
func (l *List[T]) Map(f func(T) T) *List[T] {
	if l.IsEmpty() {
		return l
	}

	// Build result in reverse
	var result *List[T] = EmptyList[T]()
	for node := l.head; node != nil; node = node.next {
		result = result.Cons(f(node.value))
	}

	// Reverse to get correct order
	return result.Reverse()
}

// Filter returns a new list containing only elements that satisfy the predicate.
func (l *List[T]) Filter(predicate func(T) bool) *List[T] {
	if l.IsEmpty() {
		return l
	}

	// Build result in reverse
	var result *List[T] = EmptyList[T]()
	for node := l.head; node != nil; node = node.next {
		if predicate(node.value) {
			result = result.Cons(node.value)
		}
	}

	// Reverse to get correct order
	return result.Reverse()
}

// FoldLeft folds the list from left to right.
func (l *List[T]) FoldLeft(initial interface{}, f func(interface{}, T) interface{}) interface{} {
	acc := initial
	for node := l.head; node != nil; node = node.next {
		acc = f(acc, node.value)
	}
	return acc
}

// FoldRight folds the list from right to left.
func (l *List[T]) FoldRight(initial interface{}, f func(T, interface{}) interface{}) interface{} {
	if l.IsEmpty() {
		return initial
	}
	return f(l.Head(), l.Tail().FoldRight(initial, f))
}

// Reverse returns a new list with elements in reverse order.
func (l *List[T]) Reverse() *List[T] {
	if l.IsEmpty() {
		return l
	}

	result := EmptyList[T]()
	for node := l.head; node != nil; node = node.next {
		result = result.Cons(node.value)
	}
	return result
}

// ForEach applies a function to each element.
func (l *List[T]) ForEach(f func(T)) {
	for node := l.head; node != nil; node = node.next {
		f(node.value)
	}
}

// ToSlice converts the list to a slice.
func (l *List[T]) ToSlice() []T {
	if l.IsEmpty() {
		return []T{}
	}

	result := make([]T, l.size)
	i := 0
	for node := l.head; node != nil; node = node.next {
		result[i] = node.value
		i++
	}
	return result
}

// String returns a string representation of the list.
func (l *List[T]) String() string {
	var sb strings.Builder
	sb.WriteString("List[")
	first := true
	for node := l.head; node != nil; node = node.next {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", node.value))
		first = false
	}
	sb.WriteString("]")
	return sb.String()
}

// Vector is a persistent immutable vector (array-like structure).
// It uses a balanced tree structure for efficient updates.
type Vector[T any] struct {
	root   *vectorNode[T]
	tail   []T
	length int
	shift  uint
}

type vectorNode[T any] struct {
	children []interface{} // Can be *vectorNode[T] or T
}

const (
	vectorNodeSize = 32
	vectorShift    = 5 // 2^5 = 32
)

// EmptyVector creates an empty vector.
func EmptyVector[T any]() *Vector[T] {
	return &Vector[T]{
		root:   nil,
		tail:   make([]T, 0, vectorNodeSize),
		length: 0,
		shift:  vectorShift,
	}
}

// VectorOf creates a vector from the given values.
func VectorOf[T any](values ...T) *Vector[T] {
	v := EmptyVector[T]()
	for _, value := range values {
		v = v.Append(value)
	}
	return v
}

// Append adds an element to the end of the vector.
// Returns a new vector with the element added.
func (v *Vector[T]) Append(value T) *Vector[T] {
	if len(v.tail) < vectorNodeSize {
		// Room in tail
		newTail := make([]T, len(v.tail)+1, vectorNodeSize)
		copy(newTail, v.tail)
		newTail[len(v.tail)] = value
		return &Vector[T]{
			root:   v.root,
			tail:   newTail,
			length: v.length + 1,
			shift:  v.shift,
		}
	}

	// Tail is full, need to push it into the tree
	newRoot := v.pushTail(v.shift, v.root, v.tail)
	newTail := []T{value}
	return &Vector[T]{
		root:   newRoot,
		tail:   newTail,
		length: v.length + 1,
		shift:  v.shift,
	}
}

func (v *Vector[T]) pushTail(level uint, node *vectorNode[T], tail []T) *vectorNode[T] {
	if node == nil {
		// Create new root node
		return &vectorNode[T]{
			children: []interface{}{tail},
		}
	}

	if level == 0 {
		// Leaf node
		children := make([]interface{}, len(node.children)+1)
		copy(children, node.children)
		children[len(node.children)] = tail
		return &vectorNode[T]{
			children: children,
		}
	}

	// Internal node
	subIdx := ((v.length - 1) >> level) & (vectorNodeSize - 1)
	child := v.pushTail(level-vectorShift, node.children[subIdx].(*vectorNode[T]), tail)
	children := make([]interface{}, len(node.children))
	copy(children, node.children)
	children[subIdx] = child
	return &vectorNode[T]{
		children: children,
	}
}

// Get returns the element at the given index.
// Panics if index is out of bounds.
func (v *Vector[T]) Get(index int) T {
	if index < 0 || index >= v.length {
		panic(fmt.Sprintf("Vector.Get: index %d out of bounds [0, %d)", index, v.length))
	}

	if index >= v.length-len(v.tail) {
		// In tail
		return v.tail[index-(v.length-len(v.tail))]
	}

	// In tree
	node := v.root
	for level := v.shift; level > 0; level -= vectorShift {
		subIdx := (index >> level) & (vectorNodeSize - 1)
		node = node.children[subIdx].(*vectorNode[T])
	}
	subIdx := index & (vectorNodeSize - 1)
	return node.children[subIdx].(T)
}

// Set replaces the element at the given index.
// Returns a new vector with the element replaced.
func (v *Vector[T]) Set(index int, value T) *Vector[T] {
	if index < 0 || index >= v.length {
		panic(fmt.Sprintf("Vector.Set: index %d out of bounds [0, %d)", index, v.length))
	}

	if index >= v.length-len(v.tail) {
		// In tail
		newTail := make([]T, len(v.tail))
		copy(newTail, v.tail)
		newTail[index-(v.length-len(v.tail))] = value
		return &Vector[T]{
			root:   v.root,
			tail:   newTail,
			length: v.length,
			shift:  v.shift,
		}
	}

	// In tree
	newRoot := v.setNode(v.shift, v.root, index, value)
	return &Vector[T]{
		root:   newRoot,
		tail:   v.tail,
		length: v.length,
		shift:  v.shift,
	}
}

func (v *Vector[T]) setNode(level uint, node *vectorNode[T], index int, value T) *vectorNode[T] {
	if level == 0 {
		// Leaf node
		children := make([]interface{}, len(node.children))
		copy(children, node.children)
		children[index&(vectorNodeSize-1)] = value
		return &vectorNode[T]{
			children: children,
		}
	}

	// Internal node
	subIdx := (index >> level) & (vectorNodeSize - 1)
	child := v.setNode(level-vectorShift, node.children[subIdx].(*vectorNode[T]), index, value)
	children := make([]interface{}, len(node.children))
	copy(children, node.children)
	children[subIdx] = child
	return &vectorNode[T]{
		children: children,
	}
}

// Length returns the number of elements in the vector.
func (v *Vector[T]) Length() int {
	return v.length
}

// IsEmpty returns true if the vector is empty.
func (v *Vector[T]) IsEmpty() bool {
	return v.length == 0
}

// Map applies a function to each element and returns a new vector.
func (v *Vector[T]) Map(f func(T) T) *Vector[T] {
	if v.IsEmpty() {
		return v
	}

	result := EmptyVector[T]()
	for i := 0; i < v.length; i++ {
		result = result.Append(f(v.Get(i)))
	}
	return result
}

// Filter returns a new vector containing only elements that satisfy the predicate.
func (v *Vector[T]) Filter(predicate func(T) bool) *Vector[T] {
	if v.IsEmpty() {
		return v
	}

	result := EmptyVector[T]()
	for i := 0; i < v.length; i++ {
		value := v.Get(i)
		if predicate(value) {
			result = result.Append(value)
		}
	}
	return result
}

// ForEach applies a function to each element.
func (v *Vector[T]) ForEach(f func(T)) {
	for i := 0; i < v.length; i++ {
		f(v.Get(i))
	}
}

// ToSlice converts the vector to a slice.
func (v *Vector[T]) ToSlice() []T {
	result := make([]T, v.length)
	for i := 0; i < v.length; i++ {
		result[i] = v.Get(i)
	}
	return result
}

// String returns a string representation of the vector.
func (v *Vector[T]) String() string {
	var sb strings.Builder
	sb.WriteString("Vector[")
	for i := 0; i < v.length; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", v.Get(i)))
	}
	sb.WriteString("]")
	return sb.String()
}

// Map is a persistent immutable hash map.
// This is a simplified implementation using a slice of key-value pairs.
// For production use, consider a more efficient data structure.
type Map[K comparable, V any] struct {
	pairs []Pair[K, V]
}

// EmptyMap creates an empty map.
func EmptyMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{pairs: []Pair[K, V]{}}
}

// MapOf creates a map from key-value pairs.
func MapOf[K comparable, V any](pairs ...Pair[K, V]) *Map[K, V] {
	m := EmptyMap[K, V]()
	for _, pair := range pairs {
		m = m.Set(pair.Key, pair.Value)
	}
	return m
}

// Pair represents a key-value pair.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// PairOf creates a key-value pair.
func PairOf[K comparable, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// Set adds or updates a key-value pair.
// Returns a new map with the pair added/updated.
func (m *Map[K, V]) Set(key K, value V) *Map[K, V] {
	// Create a new slice
	newPairs := make([]Pair[K, V], 0, len(m.pairs)+1)
	found := false

	// Copy existing pairs, updating if key exists
	for _, pair := range m.pairs {
		if pair.Key == key {
			// Update existing key
			newPairs = append(newPairs, Pair[K, V]{Key: key, Value: value})
			found = true
		} else {
			newPairs = append(newPairs, pair)
		}
	}

	// Add new key if not found
	if !found {
		newPairs = append(newPairs, Pair[K, V]{Key: key, Value: value})
	}

	return &Map[K, V]{pairs: newPairs}
}

// Get returns the value for the given key.
// Returns false as second return value if key not found.
func (m *Map[K, V]) Get(key K) (V, bool) {
	for _, pair := range m.pairs {
		if pair.Key == key {
			return pair.Value, true
		}
	}
	var zero V
	return zero, false
}

// Delete removes a key from the map.
// Returns a new map without the key.
func (m *Map[K, V]) Delete(key K) *Map[K, V] {
	// Create a new slice
	newPairs := make([]Pair[K, V], 0, len(m.pairs))

	// Copy all pairs except the one to delete
	for _, pair := range m.pairs {
		if pair.Key != key {
			newPairs = append(newPairs, pair)
		}
	}

	return &Map[K, V]{pairs: newPairs}
}

// Size returns the number of key-value pairs in the map.
func (m *Map[K, V]) Size() int {
	return len(m.pairs)
}

// IsEmpty returns true if the map is empty.
func (m *Map[K, V]) IsEmpty() bool {
	return len(m.pairs) == 0
}

// Contains returns true if the map contains the key.
func (m *Map[K, V]) Contains(key K) bool {
	_, found := m.Get(key)
	return found
}

// ForEach applies a function to each key-value pair.
func (m *Map[K, V]) ForEach(f func(K, V)) {
	for _, pair := range m.pairs {
		f(pair.Key, pair.Value)
	}
}

// Map applies a function to each value and returns a new map.
func (m *Map[K, V]) Map(f func(V) V) *Map[K, V] {
	result := EmptyMap[K, V]()
	for _, pair := range m.pairs {
		result = result.Set(pair.Key, f(pair.Value))
	}
	return result
}

// Filter returns a new map containing only key-value pairs that satisfy the predicate.
func (m *Map[K, V]) Filter(predicate func(K, V) bool) *Map[K, V] {
	result := EmptyMap[K, V]()
	for _, pair := range m.pairs {
		if predicate(pair.Key, pair.Value) {
			result = result.Set(pair.Key, pair.Value)
		}
	}
	return result
}

// Keys returns a slice of all keys in the map.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, len(m.pairs))
	for i, pair := range m.pairs {
		keys[i] = pair.Key
	}
	return keys
}

// Values returns a slice of all values in the map.
func (m *Map[K, V]) Values() []V {
	values := make([]V, len(m.pairs))
	for i, pair := range m.pairs {
		values[i] = pair.Value
	}
	return values
}

// ToSlice converts the map to a slice of key-value pairs.
func (m *Map[K, V]) ToSlice() []Pair[K, V] {
	return m.pairs
}

// String returns a string representation of the map.
func (m *Map[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("Map{")
	for i, pair := range m.pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v: %v", pair.Key, pair.Value))
	}
	sb.WriteString("}")
	return sb.String()
}

// Set is a persistent immutable set.
type Set[T comparable] struct {
	inner *Map[T, struct{}]
}

// EmptySet creates an empty set.
func EmptySet[T comparable]() *Set[T] {
	return &Set[T]{inner: EmptyMap[T, struct{}]()}
}

// SetOf creates a set from the given values.
func SetOf[T comparable](values ...T) *Set[T] {
	s := EmptySet[T]()
	for _, value := range values {
		s = s.Add(value)
	}
	return s
}

// Add adds an element to the set.
// Returns a new set with the element added.
func (s *Set[T]) Add(value T) *Set[T] {
	return &Set[T]{inner: s.inner.Set(value, struct{}{})}
}

// Remove removes an element from the set.
// Returns a new set without the element.
func (s *Set[T]) Remove(value T) *Set[T] {
	return &Set[T]{inner: s.inner.Delete(value)}
}

// Contains returns true if the set contains the element.
func (s *Set[T]) Contains(value T) bool {
	_, found := s.inner.Get(value)
	return found
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return s.inner.Size()
}

// IsEmpty returns true if the set is empty.
func (s *Set[T]) IsEmpty() bool {
	return s.inner.IsEmpty()
}

// Union returns a new set containing all elements from both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := s
	other.inner.ForEach(func(key T, _ struct{}) {
		result = result.Add(key)
	})
	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := EmptySet[T]()
	s.inner.ForEach(func(key T, _ struct{}) {
		if other.Contains(key) {
			result = result.Add(key)
		}
	})
	return result
}

// Difference returns a new set containing elements in this set but not in the other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := EmptySet[T]()
	s.inner.ForEach(func(key T, _ struct{}) {
		if !other.Contains(key) {
			result = result.Add(key)
		}
	})
	return result
}

// ForEach applies a function to each element.
func (s *Set[T]) ForEach(f func(T)) {
	s.inner.ForEach(func(key T, _ struct{}) {
		f(key)
	})
}

// ToSlice converts the set to a slice.
func (s *Set[T]) ToSlice() []T {
	return s.inner.Keys()
}

// String returns a string representation of the set.
func (s *Set[T]) String() string {
	var sb strings.Builder
	sb.WriteString("Set{")
	first := true
	s.ForEach(func(value T) {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", value))
		first = false
	})
	sb.WriteString("}")
	return sb.String()
}
