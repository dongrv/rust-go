# RustGo - Rust-like Programming in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/dongrv/rust-go.svg)](https://pkg.go.dev/github.com/dongrv/rust-go)
[![License: Apache 2.0](https://img.shields.io/badge/apache-2.0-green.svg)](https://opensource.org/licenses/apache-2-0)

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
import (
    "fmt"
    "github.com/dongrv/rust-go"
)

// Safe null handling
value := rust.Some(42)
result := rust.MapOption(value, func(x int) string {
    return fmt.Sprintf("Number: %d", x)
})

fmt.Println(result.UnwrapOr("default")) // "Number: 42"

// Chaining operations
finalResult := rust.AndThenOption(value, func(x int) rust.Option[string] {
    if x > 10 {
        return rust.Some(fmt.Sprintf("Large: %d", x))
    }
    return rust.None[string]()
})
```

### Result Type (Error Handling)
```go
import (
    "fmt"
    "github.com/dongrv/rust-go"
)

func divide(a, b int) rust.Result[int, string] {
    if b == 0 {
        return rust.Err[int, string]("division by zero")
    }
    return rust.Ok[int, string](a / b)
}

// Railway-oriented programming
computation := rust.AndThenResult(divide(10, 2), func(x int) rust.Result[int, string] {
    return rust.AndThenResult(divide(x, 2), func(y int) rust.Result[int, string] {
        return rust.Ok[int, string](y * 3)
    })
})

fmt.Println(computation) // Ok(6)
```

### Iterator Usage
```go
import (
    "fmt"
    "github.com/dongrv/rust-go"
)

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Lazy, chainable operations
result := rust.Collect(
    rust.Take(
        rust.Map(
            rust.Filter(rust.Iter(numbers), func(x int) bool { return x%2 == 0 }),
            func(x int) int { return x * 2 },
        ),
        3,
    ),
)

fmt.Println(result) // [4 8 12]

// Functional operations
sum := rust.Fold(rust.Iter(numbers), 0, func(acc, x int) int {
    return acc + x
})
fmt.Printf("Sum: %d\n", sum) // Sum: 55
```

### Chainable Collections
```go
import (
    "fmt"
    "github.com/dongrv/rust-go"
)

data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Functional pipeline
result := rust.From(data).
    Filter(func(x int) bool { return x%2 == 0 }). // Even numbers
    Map(func(x int) int { return x * 3 }).        // Triple them
    Skip(1).                                      // Skip first
    Take(3).                                      // Take next 3
    Reverse().                                    // Reverse order
    Collect()

fmt.Println(result) // [24 18 12]

// More operations
unique := rust.From([]int{1, 2, 2, 3, 3, 3}).Unique().Collect()
fmt.Println(unique) // [1 2 3]

// Partition
trueVals, falseVals := rust.From(data).Partition(func(x int) bool {
    return x > 5
})
```

## 📚 Core Concepts

### Option Pattern
The `Option[T]` type eliminates null pointer exceptions by making the presence or absence of a value explicit:

```go
// Instead of: var value *int
value := rust.Some(42)  // Explicitly has a value
empty := rust.None[int]() // Explicitly has no value

// Safe operations
result := value.UnwrapOr(0)     // 42
emptyResult := empty.UnwrapOr(0) // 0
```

### Result Pattern
The `Result[T, E]` type makes error handling explicit and type-safe:

```go
// Instead of: func() (value int, err error)
func parseInput(s string) rust.Result[int, string] {
    // Success case
    return rust.Ok[int, string](42)
    
    // Error case  
    return rust.Err[int, string]("invalid input")
}

// Railway-oriented programming
finalResult := rust.AndThenResult(
    parseInput("10"),
    func(x int) rust.Result[int, string] {
        return rust.MapResult(
            parseInput("20"),
            func(y int) int {
                return x + y
            },
        )
    },
)
```

### Iterator Pattern
Lazy, chainable iterators for efficient data processing:

```go
// Lazy evaluation - no computation until consumed
pipeline := rust.Map(
    rust.Filter(rust.Iter(data), predicate),
    transformer,
)

// Only now does computation happen
result := rust.Collect(pipeline)
```

## 🏗️ Architecture

### Package Structure
```
rust-go/
├── option.go      # Option[T] type and operations
├── result.go      # Result[T, E] type and operations  
├── iterator.go    # Iterator[T] interface and implementations
├── chainable.go   # Chainable[T] collections
├── core_test.go   # Comprehensive tests
├── go.mod         # Go module definition
├── LICENSE        # Apache 2.0 License
└── examples/      # Usage examples
    ├── option_example.go
    ├── result_example.go
    ├── iterator_example.go
    └── comprehensive_example.go
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

func (it *CustomIterator) Next() rust.Option[int] {
    if it.count >= it.max {
        return rust.None[int]()
    }
    value := it.count
    it.count++
    return rust.Some(value)
}

// Use with existing operations
iter := &CustomIterator{count: 0, max: 5}
squared := rust.Collect(rust.Map(iter, func(x int) int {
    return x * x
}))
```

### Error Recovery
```go
func processFile(path string) rust.Result[string, string] {
    content, err := os.ReadFile(path)
    if err != nil {
        return rust.Err[string, string](err.Error())
    }
    
    return rust.AndThenResult(
        validateContent(string(content)),
        func(validated string) rust.Result[string, string] {
            return rust.Ok[string, string](processContent(validated))
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

## 📖 Examples

Check out the comprehensive examples in the `examples/` directory:

- `option_example.go` - Complete Option type usage
- `result_example.go` - Result type and error handling patterns
- `iterator_example.go` - Iterator and Chainable operations
- `comprehensive_example.go` - Real-world order processing system

Run all examples:
```bash
cd examples
go run run_all.bat  # Windows
./run_all.sh        # Linux/Unix
```

Or run individual examples:
```bash
cd examples
go run option_example.go
go run result_example.go
go run iterator_example.go
go run comprehensive_example.go
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

Apache 2.0 License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Inspired by:
- Rust programming language
- Scala's functional collections  
- Haskell's monadic error handling
- Modern C++ ranges and optional types

---

**RustGo**: Write Go with Rust's elegance, enjoy Go's simplicity. 🦀 + 🐹 = ❤️