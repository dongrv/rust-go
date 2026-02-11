// package rust provides Rust-like programming constructs for Go
package rust

// Iterator is the trait for Rust-like iterators
type Iterator[T any] interface {
	// Next returns the next element in the iterator
	Next() Option[T]
}

// Pair represents a tuple of two values
type Pair[A, B any] struct {
	First  A
	Second B
}

// SliceIterator implements Iterator for slices
type SliceIterator[T any] struct {
	slice []T
	index int
}

// NewSliceIterator creates a new iterator from a slice
func NewSliceIterator[T any](slice []T) Iterator[T] {
	return &SliceIterator[T]{
		slice: slice,
		index: 0,
	}
}

// Next returns the next element in the iterator
func (it *SliceIterator[T]) Next() Option[T] {
	if it.index < len(it.slice) {
		value := it.slice[it.index]
		it.index++
		return Some(value)
	}
	return None[T]()
}

// Iter creates an iterator from a slice
func Iter[T any](slice []T) Iterator[T] {
	return NewSliceIterator(slice)
}

// MapIterator applies a function to each element
type MapIterator[T any, U any] struct {
	source Iterator[T]
	f      func(T) U
}

// Map creates an iterator which calls a function on each element
func Map[T any, U any](source Iterator[T], f func(T) U) Iterator[U] {
	return &MapIterator[T, U]{
		source: source,
		f:      f,
	}
}

func (it *MapIterator[T, U]) Next() Option[U] {
	next := it.source.Next()
	if next.IsSome() {
		return Some(it.f(next.Unwrap()))
	}
	return None[U]()
}

// FilterIterator filters elements based on a predicate
type FilterIterator[T any] struct {
	source    Iterator[T]
	predicate func(T) bool
}

// Filter creates an iterator which filters elements
func Filter[T any](source Iterator[T], predicate func(T) bool) Iterator[T] {
	return &FilterIterator[T]{
		source:    source,
		predicate: predicate,
	}
}

func (it *FilterIterator[T]) Next() Option[T] {
	for {
		next := it.source.Next()
		if next.IsNone() {
			return None[T]()
		}
		value := next.Unwrap()
		if it.predicate(value) {
			return Some(value)
		}
	}
}

// TakeIterator takes only the first n elements
type TakeIterator[T any] struct {
	source Iterator[T]
	n      int
	taken  int
}

// Take creates an iterator that yields the first n elements
func Take[T any](source Iterator[T], n int) Iterator[T] {
	return &TakeIterator[T]{
		source: source,
		n:      n,
		taken:  0,
	}
}

func (it *TakeIterator[T]) Next() Option[T] {
	if it.taken >= it.n {
		return None[T]()
	}
	next := it.source.Next()
	if next.IsSome() {
		it.taken++
	}
	return next
}

// SkipIterator skips the first n elements
type SkipIterator[T any] struct {
	source  Iterator[T]
	n       int
	skipped bool
}

// Skip creates an iterator that skips the first n elements
func Skip[T any](source Iterator[T], n int) Iterator[T] {
	return &SkipIterator[T]{
		source:  source,
		n:       n,
		skipped: false,
	}
}

func (it *SkipIterator[T]) Next() Option[T] {
	if !it.skipped {
		for i := 0; i < it.n; i++ {
			if it.source.Next().IsNone() {
				it.skipped = true
				return None[T]()
			}
		}
		it.skipped = true
	}
	return it.source.Next()
}

// ChainIterator concatenates two iterators
type ChainIterator[T any] struct {
	first       Iterator[T]
	second      Iterator[T]
	usingSecond bool
}

// Chain concatenates two iterators
func Chain[T any](first, second Iterator[T]) Iterator[T] {
	return &ChainIterator[T]{
		first:       first,
		second:      second,
		usingSecond: false,
	}
}

func (it *ChainIterator[T]) Next() Option[T] {
	if !it.usingSecond {
		next := it.first.Next()
		if next.IsSome() {
			return next
		}
		it.usingSecond = true
	}
	return it.second.Next()
}

// ZipIterator zips two iterators together
type ZipIterator[T any, U any] struct {
	first  Iterator[T]
	second Iterator[U]
}

// Zip 'zips up' two iterators into a single iterator of pairs
func Zip[T any, U any](first Iterator[T], second Iterator[U]) Iterator[Pair[T, U]] {
	return &ZipIterator[T, U]{
		first:  first,
		second: second,
	}
}

func (it *ZipIterator[T, U]) Next() Option[Pair[T, U]] {
	firstNext := it.first.Next()
	secondNext := it.second.Next()
	if firstNext.IsSome() && secondNext.IsSome() {
		return Some(Pair[T, U]{
			First:  firstNext.Unwrap(),
			Second: secondNext.Unwrap(),
		})
	}
	return None[Pair[T, U]]()
}

// EnumerateIterator adds indices to elements
type EnumerateIterator[T any] struct {
	source Iterator[T]
	index  int
}

// Enumerate creates an iterator which gives the current iteration count as well as the next value
func Enumerate[T any](source Iterator[T]) Iterator[Pair[int, T]] {
	return &EnumerateIterator[T]{
		source: source,
		index:  0,
	}
}

func (it *EnumerateIterator[T]) Next() Option[Pair[int, T]] {
	next := it.source.Next()
	if next.IsSome() {
		pair := Pair[int, T]{
			First:  it.index,
			Second: next.Unwrap(),
		}
		it.index++
		return Some(pair)
	}
	return None[Pair[int, T]]()
}

// Collect collects all elements from an iterator into a slice
func Collect[T any](iter Iterator[T]) []T {
	var result []T
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		result = append(result, next.Unwrap())
	}
	return result
}

// ForEach calls a function for each element in the iterator
func ForEach[T any](iter Iterator[T], f func(T)) {
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		f(next.Unwrap())
	}
}

// Fold folds every element into an accumulator
func Fold[T any, U any](iter Iterator[T], initial U, f func(U, T) U) U {
	acc := initial
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		acc = f(acc, next.Unwrap())
	}
	return acc
}

// Reduce reduces the elements to a single value
func Reduce[T any](iter Iterator[T], f func(T, T) T) Option[T] {
	first := iter.Next()
	if first.IsNone() {
		return None[T]()
	}

	acc := first.Unwrap()
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		acc = f(acc, next.Unwrap())
	}

	return Some(acc)
}

// All tests if every element matches a predicate
func All[T any](iter Iterator[T], predicate func(T) bool) bool {
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		if !predicate(next.Unwrap()) {
			return false
		}
	}
	return true
}

// Any tests if any element matches a predicate
func Any[T any](iter Iterator[T], predicate func(T) bool) bool {
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		if predicate(next.Unwrap()) {
			return true
		}
	}
	return false
}

// Find searches for an element that satisfies a predicate
func Find[T any](iter Iterator[T], predicate func(T) bool) Option[T] {
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		value := next.Unwrap()
		if predicate(value) {
			return Some(value)
		}
	}
	return None[T]()
}

// Count counts the number of elements in the iterator
func Count[T any](iter Iterator[T]) int {
	count := 0
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		count++
	}
	return count
}

// Last returns the last element of the iterator
func Last[T any](iter Iterator[T]) Option[T] {
	var last Option[T]
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		last = next
	}
	return last
}

// Range creates an iterator over a range of integers
func Range(start, end, step int) Iterator[int] {
	return &RangeIterator{
		current: start,
		end:     end,
		step:    step,
	}
}

type RangeIterator struct {
	current int
	end     int
	step    int
}

func (it *RangeIterator) Next() Option[int] {
	if (it.step > 0 && it.current >= it.end) || (it.step < 0 && it.current <= it.end) {
		return None[int]()
	}
	value := it.current
	it.current += it.step
	return Some(value)
}

// Once creates an iterator that yields an element exactly once
func Once[T any](value T) Iterator[T] {
	return &OnceIterator[T]{
		value:   value,
		yielded: false,
	}
}

type OnceIterator[T any] struct {
	value   T
	yielded bool
}

func (it *OnceIterator[T]) Next() Option[T] {
	if !it.yielded {
		it.yielded = true
		return Some(it.value)
	}
	return None[T]()
}

// Repeat creates an iterator that repeats an element endlessly
func Repeat[T any](value T) Iterator[T] {
	return &RepeatIterator[T]{value: value}
}

type RepeatIterator[T any] struct {
	value T
}

func (it *RepeatIterator[T]) Next() Option[T] {
	return Some(it.value)
}

// Empty creates an empty iterator
func Empty[T any]() Iterator[T] {
	return &EmptyIterator[T]{}
}

type EmptyIterator[T any] struct{}

func (it *EmptyIterator[T]) Next() Option[T] {
	return None[T]()
}
