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

### 🛡️ **Enhanced Error Handling**
- **Error**: Enhanced error type with context, stack traces, and chaining
- **Result[T]**: Railway-oriented programming with `Map`, `AndThen`, `OrElse`
- **ErrorHandler**: Fluent interface for error handling pipelines
- **ErrorChain**: Structured error tracing and flattening

### 🏗️ **Immutable Data Structures**
- **List[T]**: Persistent immutable singly-linked list
- **Vector[T]**: Persistent immutable vector with efficient updates
- **Map[K, V]**: Persistent immutable hash map
- **Set[T]**: Persistent immutable set with set operations

### 🔧 **Trait System**
- **Trait Registry**: Compile-time polymorphism with type registration
- **Dynamic Dispatch**: Runtime polymorphism through trait objects
- **Trait Composition**: Combine multiple traits for complex behaviors
- **Automatic Derivation**: Auto-generate trait implementations

### 🔄 **Functional Programming**
- Chainable operations: `map`, `filter`, `fold`, `reduce`, `take`, `skip`
- Lazy evaluation with iterators
- Railway-oriented programming for error handling
- Pattern matching inspired operations
- Immutable data structures for pure functional programming
- Trait-based polymorphism for flexible code organization

### 🛡️ **Safety & Expressiveness**
- Eliminate null pointer exceptions with `Option`
- Make error handling explicit with `Result`
- Type-safe generics throughout
- Immutable operations by default
- Compile-time trait bounds for type safety
- Structured error handling with context and tracing

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

### Enhanced Error Handling
```go
import (
    "fmt"
    "github.com/dongrv/rust-go/errors"
)

// Enhanced error with context
err := errors.New("database connection failed").
    WithContext("host", "localhost:5432").
    WithContext("database", "users").
    WithContext("operation", "query")

// Railway-oriented programming with enhanced errors
divide := func(a, b int) errors.Result[int] {
    if b == 0 {
        return errors.Err[int](errors.New("division by zero").
            WithContext("dividend", a).
            WithContext("divisor", b))
    }
    return errors.Ok(a / b)
}

computation := divide(10, 2).
    AndThen(func(x int) errors.Result[int] {
        return divide(x, 2)
    }).
    Map(func(x int) int {
        return x * 3
    })

if computation.IsOk() {
    fmt.Printf("Result: %d\n", computation.Unwrap())
}
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

// Immutable Data Structures
import "github.com/dongrv/rust-go/immutable"

// Persistent immutable list
list := immutable.ListOf(1, 2, 3, 4, 5)
transformed := list.
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * 3 }).
    Reverse()
fmt.Printf("Transformed list: %v\n", transformed)

// Persistent immutable map
productPrices := immutable.MapOf(
    immutable.PairOf("laptop", 999.99),
    immutable.PairOf("phone", 699.99),
)
discounted := productPrices.Map(func(price float64) float64 {
    return price * 0.9 // 10% discount
})
fmt.Printf("Discounted prices: %v\n", discounted)
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

// Trait System Example
import "github.com/dongrv/rust-go/trait"

type Product struct {
    ID    string
    Name  string
    Price float64
}

product := Product{ID: "P001", Name: "Laptop", Price: 999.99}

// Derive traits automatically
trait.NewDerive(product).
    Display().
    Debug().
    Clone().
    Eq().
    Default()

// Dynamic dispatch
vtable := map[string]interface{}{
    "GetName": func(p Product) string { return p.Name },
    "GetPrice": func(p Product) float64 { return p.Price },
}

obj := trait.NewTraitObject(product, vtable)
if results, err := obj.Call("GetName"); err == nil {
    fmt.Printf("Product name: %s\n", results[0])
}
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
├── errors/        # Enhanced error handling
│   ├── errors.go      # Error, Result[T], ErrorHandler
│   └── errors_test.go # Comprehensive tests
├── immutable/     # Immutable data structures
│   ├── immutable.go   # List, Vector, Map, Set
│   └── immutable_test.go
├── trait/         # Trait system
│   ├── trait.go       # Trait registry, dynamic dispatch
│   └── trait_test.go
├── pattern/       # Pattern matching
│   ├── match.go       # Pattern matching utilities
│   └── match_test.go
└── examples/      # Usage examples
    └── examples.go    # Unified examples with CLI
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
- **Persistent data structures** with structural sharing
- **Efficient trait dispatch** with registry caching

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

// Enhanced Error Handling Pipeline
import "github.com/dongrv/rust-go/errors"

errors.Handle(nil).
    Then(func() error {
        fmt.Println("Step 1: Connecting to database...")
        return nil
    }).
    Then(func() error {
        fmt.Println("Step 2: Executing query...")
        return errors.New("query timeout")
    }).
    Map(func(e error) error {
        return errors.Wrap(e, "database operation failed")
    }).
    Log(func(format string, args ...interface{}) {
        fmt.Printf("Log: "+format+"\n", args...)
    })

// Immutable Data Structure Operations
import "github.com/dongrv/rust-go/immutable"

// Create immutable collections
list := immutable.ListOf(1, 2, 3, 4, 5)
vector := immutable.VectorOf("apple", "banana", "cherry")
mapping := immutable.MapOf(
    immutable.PairOf("key1", "value1"),
    immutable.PairOf("key2", "value2"),
)
set := immutable.SetOf(1, 2, 3, 4, 5)

// All operations return new instances
newList := list.Cons(0)           // Add to front
newVector := vector.Append("date") // Add to end
newMapping := mapping.Set("key3", "value3") // Add/update
newSet := set.Add(6)              // Add element

// Original collections remain unchanged
fmt.Printf("Original list size: %d\n", list.Size())      // 5
fmt.Printf("New list size: %d\n", newList.Size())        // 6
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
| `std::collections` | `immutable.List/Vector/Map/Set` |
| Traits | `trait` package with registry and dispatch |
| `thiserror`, `anyhow` | `errors` package with context and chaining |

## 🎯 Best Practices

1. **Use Option for nullable values** instead of nil pointers
2. **Use Result for operations that can fail** instead of multiple return values
3. **Prefer functional operations** (`map`, `filter`, `fold`) over imperative loops
4. **Chain operations** for better readability
5. **Use lazy iterators** for large datasets
6. **Leverage type inference** with generics
7. **Use immutable data structures** for thread safety and predictability
8. **Apply trait-based design** for better code organization and reuse
9. **Utilize enhanced error handling** for better debugging and context
10. **Combine features** for expressive and safe code patterns

## 📖 Examples

Check out the comprehensive examples in the `examples/` directory:

Run examples with the unified CLI:
```bash
cd examples

# Run all examples
go run examples.go all

# Run core examples (Option, Result, Iterator, Pattern)
go run examples.go core

# Run enhanced examples (Error handling, Immutable, Traits)
go run examples.go enhanced

# Run individual examples
go run examples.go option      # Option type
go run examples.go result      # Result type  
go run examples.go iterator    # Iterator and Chainable
go run examples.go pattern     # Pattern matching
go run examples.go inventory   # Combined product inventory example
```

The unified `examples.go` includes:
- **Option Type**: Null-safe programming patterns
- **Result Type**: Railway-oriented error handling
- **Iterator**: Lazy, chainable operations
- **Pattern Matching**: Expressive control flow
- **Enhanced Error Handling**: Context-aware errors with stack traces
- **Immutable Data Structures**: Persistent List, Vector, Map, Set
- **Trait System**: Compile-time and runtime polymorphism
- **Product Inventory**: Combined example using all features

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
- Clojure's persistent data structures
- Swift's protocol-oriented programming

---

**RustGo**: Write Go with Rust's elegance, enjoy Go's simplicity. 🦀 + 🐹 = ❤️
