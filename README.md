# RustGo - Rust-like Programming in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/dongrv/rust-go.svg)](https://pkg.go.dev/github.com/dongrv/rust-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/license/apache-2-0)

RustGo is a comprehensive framework that brings Rust's elegant programming patterns to Go. Write Go code with Rust's expressiveness, safety, and functional programming capabilities while maintaining Go's simplicity and performance.

## 🚀 Features

### 🎯 **Core Types**
- **Option[T]**: Null-safe optional values with `Some` and `None` variants
- **Result[T, E]**: Type-safe error handling with `Ok` and `Err` variants
- **Iterator[T]**: Lazy, chainable iterators with Rust-like API
- **Chainable[T]**: Functional operations on slices and collections

### 🔄 **Functional Programming**
- Chainable operations: `map`, `filter`, `fold`, `reduce`, `take`, `skip`
- Lazy evaluation with iterators
- Railway-oriented programming for error handling
- Pattern matching inspired operations

### 🛡️ **Safety & Expressiveness**
- Eliminate null pointer exceptions with `Option`
- Make error handling explicit with `Result`
- Type-safe generics throughout
- Immutable operations by default

## 📦 Installation

```bash
go get github.com/dongrv/rust-go
```

## 🎮 Quick Start

### Option Type (Null Safety)
```go
import "github.com/dongrv/rust-go/pkg/rustgo/core"

// Safe null handling
value := core.Some(42)
result := core.MapOption(value, func(x int) string {
    return fmt.Sprintf("Number: %d", x)
})

fmt.Println(result.UnwrapOr("default")) // "Number: 42"

// Chaining operations
finalResult := core.AndThenOption(value, func(x int) core.Option[string] {
    if x > 10 {
        return core.Some(fmt.Sprintf("Large: %d", x))
    }
    return core.None[string]()
})
```

### Result Type (Error Handling)
```go
func divide(a, b int) core.Result[int, string] {
    if b == 0 {
        return core.Err[int, string]("division by zero")
    }
    return core.Ok[int, string](a / b)
}

// Railway-oriented programming
computation := core.AndThenResult(divide(10, 2), func(x int) core.Result[int, string] {
    return core.AndThenResult(divide(x, 2), func(y int) core.Result[int, string] {
        return core.Ok[int, string](y * 3)
    })
})

fmt.Println(computation) // Ok(6)
```

### Iterator Usage
```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Lazy, chainable operations
result := core.Collect(
    core.Take(
        core.Map(
            core.Filter(core.Iter(numbers), func(x int) bool { return x%2 == 0 }),
            func(x int) int { return x * 2 },
        ),
        3,
    ),
)

fmt.Println(result) // [4 8 12]

// Functional operations
sum := core.Fold(core.Iter(numbers), 0, func(acc, x int) int {
    return acc + x
})
fmt.Printf("Sum: %d\n", sum) // Sum: 55
```

### Chainable Collections
```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Functional pipeline
result := core.From(data).
    Filter(func(x int) bool { return x%2 == 0 }). // Even numbers
    Map(func(x int) int { return x * 3 }).        // Triple them
    Skip(1).                                      // Skip first
    Take(3).                                      // Take next 3
    Reverse().                                    // Reverse order
    Collect()

fmt.Println(result) // [24 18 12]

// More operations
unique := core.From([]int{1, 2, 2, 3, 3, 3}).Unique().Collect()
fmt.Println(unique) // [1 2 3]

// Partition
trueVals, falseVals := core.From(data).Partition(func(x int) bool {
    return x > 5
})
```

## 📚 Core Concepts

### Option Pattern
The `Option[T]` type eliminates null pointer exceptions by making the presence or absence of a value explicit:

```go
// Instead of: var value *int
value := core.Some(42)  // Explicitly has a value
empty := core.None[int]() // Explicitly has no value

// Safe operations
result := value.UnwrapOr(0)     // 42
emptyResult := empty.UnwrapOr(0) // 0
```

### Result Pattern
The `Result[T, E]` type makes error handling explicit and type-safe:

```go
// Instead of: func() (value int, err error)
func parseInput(s string) core.Result[int, string] {
    // Success case
    return core.Ok[int, string](42)
    
    // Error case  
    return core.Err[int, string]("invalid input")
}

// Railway-oriented programming
finalResult := parseInput("10").
    AndThen(func(x int) core.Result[int, string] {
        return parseInput("20").Map(func(y int) int {
            return x + y
        })
    })
```

### Iterator Pattern
Lazy, chainable iterators for efficient data processing:

```go
// Lazy evaluation - no computation until consumed
pipeline := core.Map(
    core.Filter(core.Iter(data), predicate),
    transformer,
)

// Only now does computation happen
result := core.Collect(pipeline)
```

## 🏗️ Architecture

### Package Structure
```
pkg/rustgo/core/
├── option.go      # Option[T] type and operations
├── result.go      # Result[T, E] type and operations  
├── iterator.go    # Iterator[T] interface and implementations
└── chainable.go   # Chainable[T] collections
```

### Type Safety
All types use Go generics for compile-time type safety:
- `Option[T]` works with any type `T`
- `Result[T, E]` is generic over both success and error types
- `Iterator[T]` and `Chainable[T]` provide type-safe operations

### Performance
- **Zero-cost abstractions** where possible
- **Lazy evaluation** for iterators (compute only when needed)
- **Minimal allocations** in hot paths
- **Compile-time optimizations** through generics

## 🔧 Advanced Usage

### Custom Iterator
```go
type CustomIterator struct {
    count int
    max   int
}

func (it *CustomIterator) Next() core.Option[int] {
    if it.count >= it.max {
        return core.None[int]()
    }
    value := it.count
    it.count++
    return core.Some(value)
}

// Use with existing operations
iter := &CustomIterator{count: 0, max: 5}
squared := core.Collect(core.Map(iter, func(x int) int {
    return x * x
}))
```

### Error Recovery
```go
func processFile(path string) core.Result[string, string] {
    content, err := os.ReadFile(path)
    if err != nil {
        return core.Err[string, string](err.Error())
    }
    
    return core.AndThenResult(
        validateContent(string(content)),
        func(validated string) core.Result[string, string] {
            return core.Ok[string, string](processContent(validated))
        },
    )
}

// Recover from errors
finalResult := processFile("data.txt").
    UnwrapOrElse(func(err string) string {
        return fmt.Sprintf("Error: %s, using default", err)
    })
```

## 📊 Comparison

### vs Traditional Go
| Traditional Go | RustGo |
|----------------|--------|
| `var value *int` | `Option[int]` |
| `value, err := fn()` | `Result[T, error]` |
| `for _, v := range slice` | `Iter(slice).ForEach()` |
| Manual error checking | Railway-oriented programming |

### vs Rust
| Rust | RustGo |
|------|--------|
| `Option<T>` | `Option[T]` |
| `Result<T, E>` | `Result[T, E]` |
| `iter.map().filter()` | `Map(Filter(iter))` |
| Pattern matching | Method-based matching |

## 🎯 Best Practices

1. **Use Option for nullable values** instead of nil pointers
2. **Use Result for operations that can fail** instead of multiple return values
3. **Prefer functional operations** (`map`, `filter`, `fold`) over imperative loops
4. **Chain operations** for better readability
5. **Use lazy iterators** for large datasets
6. **Leverage type inference** with generics

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Inspired by:
- Rust programming language
- Scala's functional collections  
- Haskell's monadic error handling
- Modern C++ ranges and optional types

---

**RustGo**: Write Go with Rust's elegance, enjoy Go's simplicity. 🦀 + 🐹 = ❤️
