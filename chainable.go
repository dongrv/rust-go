// package rust provides Rust-like programming constructs for Go
package rust

// Chainable provides Rust-like chainable operations for slices
type Chainable[T any] struct {
	data []T
}

type ChainablePair[A any, B any] struct {
	data []Pair[A, B]
}

// ChainableSlice is a wrapper for slices of slices to avoid instantiation cycles
type ChainableSlice[T any] struct {
	data [][]T
}

// NewChainable creates a new Chainable from a slice
func NewChainable[T any](data []T) *Chainable[T] {
	return &Chainable[T]{data: data}
}

// From creates a Chainable from a slice
func From[T any](data []T) *Chainable[T] {
	return NewChainable(data)
}

// Collect returns the underlying slice
func (c *Chainable[T]) Collect() []T {
	return c.data
}

// Iter returns an iterator over the data
func (c *Chainable[T]) Iter() Iterator[T] {
	return Iter(c.data)
}

// Map applies a function to each element
func (c *Chainable[T]) Map(f func(T) T) *Chainable[T] {
	result := make([]T, len(c.data))
	for i, v := range c.data {
		result[i] = f(v)
	}
	return NewChainable(result)
}

// Filter filters elements based on a predicate
func (c *Chainable[T]) Filter(predicate func(T) bool) *Chainable[T] {
	var result []T
	for _, v := range c.data {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return NewChainable(result)
}

// Fold folds elements into an accumulator
func (c *Chainable[T]) Fold(initial T, f func(T, T) T) T {
	acc := initial
	for _, v := range c.data {
		acc = f(acc, v)
	}
	return acc
}

// Reduce reduces elements to a single value
func (c *Chainable[T]) Reduce(f func(T, T) T) Option[T] {
	if len(c.data) == 0 {
		return None[T]()
	}
	acc := c.data[0]
	for _, v := range c.data[1:] {
		acc = f(acc, v)
	}
	return Some(acc)
}

// ForEach calls a function for each element
func (c *Chainable[T]) ForEach(f func(T)) {
	for _, v := range c.data {
		f(v)
	}
}

// All returns true if all elements satisfy the predicate
func (c *Chainable[T]) All(predicate func(T) bool) bool {
	for _, v := range c.data {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if any element satisfies the predicate
func (c *Chainable[T]) Any(predicate func(T) bool) bool {
	for _, v := range c.data {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Find returns the first element that satisfies the predicate
func (c *Chainable[T]) Find(predicate func(T) bool) Option[T] {
	for _, v := range c.data {
		if predicate(v) {
			return Some(v)
		}
	}
	return None[T]()
}

// Take takes the first n elements
func (c *Chainable[T]) Take(n int) *Chainable[T] {
	if n <= 0 {
		return NewChainable([]T{})
	}
	if n >= len(c.data) {
		return NewChainable(c.data)
	}
	return NewChainable(c.data[:n])
}

// Skip skips the first n elements
func (c *Chainable[T]) Skip(n int) *Chainable[T] {
	if n <= 0 {
		return NewChainable(c.data)
	}
	if n >= len(c.data) {
		return NewChainable([]T{})
	}
	return NewChainable(c.data[n:])
}

// Reverse reverses the order of elements
func (c *Chainable[T]) Reverse() *Chainable[T] {
	result := make([]T, len(c.data))
	for i, v := range c.data {
		result[len(c.data)-1-i] = v
	}
	return NewChainable(result)
}

// Unique returns a new Chainable with duplicate elements removed
func (c *Chainable[T]) Unique() *Chainable[T] {
	seen := make(map[any]bool)
	var result []T
	for _, v := range c.data {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return NewChainable(result)
}

// Partition partitions elements into two groups
func (c *Chainable[T]) Partition(predicate func(T) bool) (*Chainable[T], *Chainable[T]) {
	var trueElems []T
	var falseElems []T
	for _, v := range c.data {
		if predicate(v) {
			trueElems = append(trueElems, v)
		} else {
			falseElems = append(falseElems, v)
		}
	}
	return NewChainable(trueElems), NewChainable(falseElems)
}

// Zip zips with another slice
func (c *Chainable[T]) Zip(other []T) *ChainablePair[T, T] {
	minLen := len(c.data)
	if len(other) < minLen {
		minLen = len(other)
	}
	result := make([]Pair[T, T], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = Pair[T, T]{
			First:  c.data[i],
			Second: other[i],
		}
	}
	return &ChainablePair[T, T]{data: result}
}

// Enumerate adds indices to elements
func (c *Chainable[T]) Enumerate() *ChainablePair[int, T] {
	result := make([]Pair[int, T], len(c.data))
	for i, v := range c.data {
		result[i] = Pair[int, T]{
			First:  i,
			Second: v,
		}
	}
	return &ChainablePair[int, T]{data: result}
}

// FlatMap maps each element to a slice and flattens the result
func (c *Chainable[T]) FlatMap(f func(T) []T) *Chainable[T] {
	var result []T
	for _, v := range c.data {
		result = append(result, f(v)...)
	}
	return NewChainable(result)
}

// Chunk splits the data into chunks of specified size
func (c *Chainable[T]) Chunk(size int) *ChainableSlice[T] {
	if size <= 0 {
		return &ChainableSlice[T]{data: [][]T{}}
	}
	var result [][]T
	for i := 0; i < len(c.data); i += size {
		end := i + size
		if end > len(c.data) {
			end = len(c.data)
		}
		result = append(result, c.data[i:end])
	}
	return &ChainableSlice[T]{data: result}
}

// Window creates sliding windows of specified size
func (c *Chainable[T]) Window(size int) *ChainableSlice[T] {
	if size <= 0 || size > len(c.data) {
		return &ChainableSlice[T]{data: [][]T{}}
	}
	var result [][]T
	for i := 0; i <= len(c.data)-size; i++ {
		result = append(result, c.data[i:i+size])
	}
	return &ChainableSlice[T]{data: result}
}

// Append appends elements
func (c *Chainable[T]) Append(elements ...T) *Chainable[T] {
	result := make([]T, len(c.data)+len(elements))
	copy(result, c.data)
	copy(result[len(c.data):], elements)
	return NewChainable(result)
}

// Prepend prepends elements
func (c *Chainable[T]) Prepend(elements ...T) *Chainable[T] {
	result := make([]T, len(elements)+len(c.data))
	copy(result, elements)
	copy(result[len(elements):], c.data)
	return NewChainable(result)
}

// Concat concatenates multiple chainables
func (c *Chainable[T]) Concat(others ...*Chainable[T]) *Chainable[T] {
	totalLen := len(c.data)
	for _, other := range others {
		totalLen += len(other.data)
	}
	result := make([]T, totalLen)
	copy(result, c.data)
	offset := len(c.data)
	for _, other := range others {
		copy(result[offset:], other.data)
		offset += len(other.data)
	}
	return NewChainable(result)
}

// Helper functions

// Of creates a chainable from variadic arguments
func Of[T any](elements ...T) *Chainable[T] {
	return NewChainable(elements)
}

// Range creates a chainable of integers in a range
func RangeChainable(start, end, step int) *Chainable[int] {
	if step == 0 {
		return NewChainable([]int{})
	}
	var result []int
	if step > 0 {
		for i := start; i < end; i += step {
			result = append(result, i)
		}
	} else {
		for i := start; i > end; i += step {
			result = append(result, i)
		}
	}
	return NewChainable(result)
}

// Generate creates a chainable by calling a generator n times
func Generate[T any](n int, generator func(int) T) *Chainable[T] {
	result := make([]T, n)
	for i := range result {
		result[i] = generator(i)
	}
	return NewChainable(result)
}

// Empty creates an empty chainable
func EmptyChainable[T any]() *Chainable[T] {
	return NewChainable([]T{})
}

// Single creates a chainable with a single element
func Single[T any](value T) *Chainable[T] {
	return NewChainable([]T{value})
}

// Collect returns the underlying slice for ChainableSlice
func (c *ChainableSlice[T]) Collect() [][]T {
	return c.data
}

// Map applies a function to each element in ChainableSlice
func (c *ChainableSlice[T]) Map(f func([]T) []T) *ChainableSlice[T] {
	result := make([][]T, len(c.data))
	for i, v := range c.data {
		result[i] = f(v)
	}
	return &ChainableSlice[T]{data: result}
}

// Collect returns the underlying slice for ChainablePair
func (c *ChainablePair[A, B]) Collect() []Pair[A, B] {
	return c.data
}

// Map applies a function to each element in ChainablePair
func (c *ChainablePair[A, B]) Map(f func(Pair[A, B]) Pair[A, B]) *ChainablePair[A, B] {
	result := make([]Pair[A, B], len(c.data))
	for i, v := range c.data {
		result[i] = f(v)
	}
	return &ChainablePair[A, B]{data: result}
}
