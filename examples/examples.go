// Package main provides comprehensive examples of Rust-like features in Go.
// This package demonstrates the usage of Option, Result, Iterator, Pattern Matching,
// and enhanced features like error handling, immutable data structures, and trait system.
package main

import (
	"fmt"
	"os"

	"github.com/dongrv/rust-go"
	"github.com/dongrv/rust-go/errors"
	"github.com/dongrv/rust-go/immutable"
	"github.com/dongrv/rust-go/pattern"
	"github.com/dongrv/rust-go/trait"
)

// RunAll runs all RustGo examples.
func RunAll() {
	fmt.Println("=== RustGo Examples Runner ===")
	fmt.Println()

	// Run core examples
	fmt.Println("Running Core Examples...")
	fmt.Println("========================")
	RunCoreExamples()
	fmt.Println()

	// Run enhanced examples
	fmt.Println("Running Enhanced Features Examples...")
	fmt.Println("=====================================")
	RunEnhancedExamples()
	fmt.Println()

	fmt.Println("=== All Examples Completed ===")
}

// RunCoreExamples runs examples of core Rust-like features.
func RunCoreExamples() {
	fmt.Println("=== Core Rust-like Features ===")
	fmt.Println()

	fmt.Println("1. Option Type Example")
	fmt.Println("=======================")
	RunOptionExample()
	fmt.Println()

	fmt.Println("2. Result Type Example")
	fmt.Println("======================")
	RunResultExample()
	fmt.Println()

	fmt.Println("3. Iterator Example")
	fmt.Println("===================")
	RunIteratorExample()
	fmt.Println()

	fmt.Println("4. Pattern Matching Example")
	fmt.Println("===========================")
	RunPatternExample()
}

// RunEnhancedExamples runs examples of enhanced features.
func RunEnhancedExamples() {
	fmt.Println("=== Enhanced Features ===")
	fmt.Println()

	fmt.Println("1. Enhanced Error Handling")
	fmt.Println("==========================")
	RunErrorHandlingExample()
	fmt.Println()

	fmt.Println("2. Immutable Data Structures")
	fmt.Println("============================")
	RunImmutableExample()
	fmt.Println()

	fmt.Println("3. Trait System")
	fmt.Println("===============")
	RunTraitExample()
	fmt.Println()

	fmt.Println("4. Combined Example: Product Inventory")
	fmt.Println("======================================")
	RunProductInventoryExample()
}

// RunOptionExample demonstrates the Option type.
func RunOptionExample() {
	fmt.Println("=== Option Type Example ===")
	fmt.Println()

	// Basic Option creation
	fmt.Println("1. Basic Option Creation:")
	someValue := rust.Some(42)
	noneValue := rust.None[int]()
	fmt.Printf("  Some(42): %v\n", someValue)
	fmt.Printf("  None[int](): %v\n", noneValue)
	fmt.Printf("  IsSome: %v, IsNone: %v\n", someValue.IsSome(), someValue.IsNone())

	// Safe unwrapping
	fmt.Println("\n2. Safe Unwrapping:")
	fmt.Printf("  UnwrapOr(Some(42), 0): %d\n", someValue.UnwrapOr(0))
	fmt.Printf("  UnwrapOr(None[int](), 0): %d\n", noneValue.UnwrapOr(0))

	// Chaining operations
	fmt.Println("\n3. Chaining Operations:")
	result := rust.AndThenOption(someValue, func(x int) rust.Option[string] {
		if x > 10 {
			return rust.Some(fmt.Sprintf("Large number: %d", x))
		}
		return rust.None[string]()
	})
	fmt.Printf("  AndThen result: %v\n", result)

	// Real-world example
	fmt.Println("\n4. Real-world Example:")
	parseInt := func(s string) rust.Option[int] {
		var n int
		_, err := fmt.Sscanf(s, "%d", &n)
		if err != nil {
			return rust.None[int]()
		}
		return rust.Some(n)
	}

	input := "123"
	parsed := parseInt(input)
	processed := rust.MapOption(parsed, func(x int) string {
		return fmt.Sprintf("Parsed: %d", x)
	})
	fmt.Printf("  Parse '%s': %v\n", input, processed)
}

// RunResultExample demonstrates the Result type.
func RunResultExample() {
	fmt.Println("=== Result Type Example ===")
	fmt.Println()

	// Basic Result creation
	fmt.Println("1. Basic Result Creation:")
	okResult := rust.Ok[int, string](42)
	errResult := rust.Err[int, string]("division by zero")
	fmt.Printf("  Ok(42): %v\n", okResult)
	fmt.Printf("  Err(\"division by zero\"): %v\n", errResult)
	fmt.Printf("  IsOk: %v, IsErr: %v\n", okResult.IsOk(), okResult.IsErr())

	// Railway-oriented programming
	fmt.Println("\n2. Railway-oriented Programming:")
	divide := func(a, b int) rust.Result[int, string] {
		if b == 0 {
			return rust.Err[int, string]("division by zero")
		}
		return rust.Ok[int, string](a / b)
	}

	computation := rust.AndThenResult(divide(10, 2), func(x int) rust.Result[int, string] {
		return rust.AndThenResult(divide(x, 2), func(y int) rust.Result[int, string] {
			return rust.Ok[int, string](y * 3)
		})
	})
	fmt.Printf("  Computation result: %v\n", computation)

	// Error handling patterns
	fmt.Println("\n3. Error Handling Patterns:")
	userInput := "not a number"
	parseResult := func(s string) rust.Result[int, string] {
		var n int
		_, err := fmt.Sscanf(s, "%d", &n)
		if err != nil {
			return rust.Err[int, string](fmt.Sprintf("parse error: %v", err))
		}
		return rust.Ok[int, string](n)
	}

	result := parseResult(userInput)
	finalResult := rust.MapResult(result, func(x int) string {
		return fmt.Sprintf("Success: %d", x)
	})
	fmt.Printf("  Parse '%s': %v\n", userInput, finalResult)
}

// RunIteratorExample demonstrates the Iterator type.
func RunIteratorExample() {
	fmt.Println("=== Iterator Example ===")
	fmt.Println()

	// Basic iterator operations
	fmt.Println("1. Basic Iterator Operations:")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iter := rust.Iter(numbers)

	// Filter even numbers, double them, take first 3
	result := rust.Collect(
		rust.Take(
			rust.Map(
				rust.Filter(iter, func(x int) bool { return x%2 == 0 }),
				func(x int) int { return x * 2 },
			),
			3,
		),
	)
	fmt.Printf("  Original: %v\n", numbers)
	fmt.Printf("  Even numbers * 2 (first 3): %v\n", result)

	// Chainable operations
	fmt.Println("\n2. Chainable Operations:")
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chainableResult := rust.From(data).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * 3 }).
		Skip(1).
		Take(2).
		Collect()
	fmt.Printf("  Chainable result: %v\n", chainableResult)

	// Functional operations
	fmt.Println("\n3. Functional Operations:")
	sum := rust.Fold(rust.Iter(numbers), 0, func(acc, x int) int {
		return acc + x
	})
	fmt.Printf("  Sum of 1..10: %d\n", sum)

	anyEven := rust.Any(rust.Iter(numbers), func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("  Any even numbers: %v\n", anyEven)

	allPositive := rust.All(rust.Iter(numbers), func(x int) bool {
		return x > 0
	})
	fmt.Printf("  All positive numbers: %v\n", allPositive)
}

// RunPatternExample demonstrates pattern matching.
func RunPatternExample() {
	fmt.Println("=== Pattern Matching Example ===")
	fmt.Println()

	// Option pattern matching
	fmt.Println("1. Option Pattern Matching:")
	optionValue := rust.Some(42)
	pattern.Match(optionValue).
		Some(func(x int) {
			fmt.Printf("  Got Some value: %d\n", x)
		}).
		None(func() {
			fmt.Println("  Got None")
		}).
		Exhaustive()

	// Result pattern matching
	fmt.Println("\n2. Result Pattern Matching:")
	resultValue := rust.Ok[int, string](42)
	pattern.Match(resultValue).
		Ok(func(x int) {
			fmt.Printf("  Success: %d\n", x)
		}).
		Err(func(err string) {
			fmt.Printf("  Error: %s\n", err)
		}).
		Exhaustive()

	// Value and type matching
	fmt.Println("\n3. Value and Type Matching:")
	var value interface{} = "hello"
	pattern.Match(value).
		Value("hello", func() {
			fmt.Println("  Exact match: hello")
		}).
		Type(func(s string) {
			fmt.Printf("  Type match: string with value '%s'\n", s)
		}).
		Default(func() {
			fmt.Println("  No match")
		})

	// String pattern matching
	fmt.Println("\n4. String Pattern Matching:")
	str := "hello_world.txt"
	pattern.MatchString(str).
		Prefix("hello", func(s string) {
			fmt.Printf("  Starts with 'hello': %s\n", s)
		}).
		Suffix(".txt", func(s string) {
			fmt.Printf("  Ends with '.txt': %s\n", s)
		}).
		Contains("world", func(s string) {
			fmt.Printf("  Contains 'world': %s\n", s)
		}).
		Default(func() {
			fmt.Println("  No string pattern match")
		})
}

// RunErrorHandlingExample demonstrates enhanced error handling.
func RunErrorHandlingExample() {
	fmt.Println("=== Enhanced Error Handling Example ===")
	fmt.Println()

	// Enhanced Error type
	fmt.Println("1. Enhanced Error Type:")
	err := errors.New("database connection failed").
		WithContext("host", "localhost:5432").
		WithContext("database", "users").
		WithContext("operation", "query")

	fmt.Printf("  Error: %v\n", err)
	fmt.Printf("  Error details:\n%s\n", err.String())

	// Result type with enhanced errors
	fmt.Println("\n2. Result Type with Enhanced Errors:")
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
		fmt.Printf("  Computation result: %d\n", computation.Unwrap())
	}

	// ErrorHandler fluent interface
	fmt.Println("\n3. ErrorHandler Fluent Interface:")
	errors.Handle(nil).
		Then(func() error {
			fmt.Println("  Step 1: Connecting to database...")
			return nil
		}).
		Then(func() error {
			fmt.Println("  Step 2: Executing query...")
			return errors.New("query timeout")
		}).
		Map(func(e error) error {
			return errors.Wrap(e, "database operation failed")
		}).
		Log(func(format string, args ...interface{}) {
			fmt.Printf("  Log: "+format+"\n", args...)
		})
}

// RunImmutableExample demonstrates immutable data structures.
func RunImmutableExample() {
	fmt.Println("=== Immutable Data Structures Example ===")
	fmt.Println()

	// Immutable List
	fmt.Println("1. Immutable List:")
	list := immutable.ListOf(1, 2, 3, 4, 5)
	fmt.Printf("  Original list: %v\n", list)

	transformed := list.
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * 3 }).
		Reverse()
	fmt.Printf("  Transformed (even*3, reversed): %v\n", transformed)

	// Immutable Map
	fmt.Println("\n2. Immutable Map:")
	productPrices := immutable.MapOf(
		immutable.PairOf("laptop", 999.99),
		immutable.PairOf("phone", 699.99),
		immutable.PairOf("tablet", 399.99),
	)
	fmt.Printf("  Original prices: %v\n", productPrices)

	discounted := productPrices.Map(func(price float64) float64 {
		return price * 0.9 // 10% discount
	})
	fmt.Printf("  After 10%% discount: %v\n", discounted)

	// Immutable Set
	fmt.Println("\n3. Immutable Set:")
	techSkills := immutable.SetOf("Go", "Rust", "Python", "JavaScript")
	backendSkills := immutable.SetOf("Go", "Rust", "Python", "Java")

	fmt.Printf("  Tech skills: %v\n", techSkills)
	fmt.Printf("  Backend skills: %v\n", backendSkills)
	fmt.Printf("  Common skills: %v\n", techSkills.Intersection(backendSkills))
	fmt.Printf("  Unique tech skills: %v\n", techSkills.Difference(backendSkills))

	// Demonstrating immutability
	fmt.Println("\n4. Immutability Demonstration:")
	original := immutable.ListOf(1, 2, 3)
	modified := original.Cons(0)
	fmt.Printf("  Original: %v (size: %d)\n", original, original.Size())
	fmt.Printf("  Modified: %v (size: %d)\n", modified, modified.Size())
	fmt.Printf("  Original unchanged: %v\n", original)
}

// RunTraitExample demonstrates the trait system.
func RunTraitExample() {
	fmt.Println("=== Trait System Example ===")
	fmt.Println()

	// Clear registry for clean example
	trait.ClearRegistry()

	// Define a custom type
	type Product struct {
		ID    string
		Name  string
		Price float64
	}

	product := Product{ID: "P001", Name: "Laptop", Price: 999.99}

	// Derive traits automatically
	fmt.Println("1. Automatic Trait Derivation:")
	trait.NewDerive(product).
		Display().
		Debug().
		Clone().
		Eq().
		Default()

	fmt.Printf("  Has Display trait: %v\n", trait.HasTrait("Display", product))
	fmt.Printf("  Has Debug trait: %v\n", trait.HasTrait("Debug", product))
	fmt.Printf("  Has Clone trait: %v\n", trait.HasTrait("Clone", product))

	// Dynamic dispatch
	fmt.Println("\n2. Dynamic Dispatch:")
	vtable := map[string]interface{}{
		"GetID":    func(p Product) string { return p.ID },
		"GetName":  func(p Product) string { return p.Name },
		"GetPrice": func(p Product) float64 { return p.Price },
		"ApplyDiscount": func(p Product, percent float64) Product {
			return Product{
				ID:    p.ID,
				Name:  p.Name,
				Price: p.Price * (1 - percent/100),
			}
		},
	}

	obj := trait.NewTraitObject(product, vtable)

	if results, err := obj.Call("GetName"); err == nil {
		fmt.Printf("  Product name: %s\n", results[0])
	}

	if results, err := obj.Call("ApplyDiscount", 10.0); err == nil {
		discounted := results[0].(Product)
		fmt.Printf("  Price after 10%% discount: $%.2f\n", discounted.Price)
	}

	// Trait composition
	fmt.Println("\n3. Trait Composition:")
	comp := trait.Compose("Display", "Debug", "Clone")
	impl := comp.Implement(product)

	if _, found := impl.GetTrait("Display"); found {
		fmt.Println("  Product implements Display trait")
	}
	if _, found := impl.GetTrait("Debug"); found {
		fmt.Println("  Product implements Debug trait")
	}
}

// RunProductInventoryExample demonstrates a combined example using all features.
func RunProductInventoryExample() {
	fmt.Println("=== Combined Example: Product Inventory ===")
	fmt.Println()

	// Define Product type
	type Product struct {
		ID     string
		Name   string
		Price  float64
		Stock  int
		Active bool
	}

	// Derive traits for Product
	trait.NewDerive(Product{}).
		Display().
		Eq()

	// Create immutable product inventory
	inventory := immutable.MapOf(
		immutable.PairOf("P001", Product{
			ID:     "P001",
			Name:   "Laptop Pro",
			Price:  1299.99,
			Stock:  25,
			Active: true,
		}),
		immutable.PairOf("P002", Product{
			ID:     "P002",
			Name:   "Smartphone X",
			Price:  899.99,
			Stock:  50,
			Active: true,
		}),
		immutable.PairOf("P003", Product{
			ID:     "P003",
			Name:   "Tablet Lite",
			Price:  399.99,
			Stock:  0,
			Active: false,
		}),
	)

	// Function to find product with error handling
	findProduct := func(id string) errors.Result[Product] {
		if product, found := inventory.Get(id); found {
			return errors.Ok(product)
		}
		return errors.Err[Product](errors.Errorf("product not found: %s", id))
	}

	// Use the inventory system
	fmt.Println("\nProduct Inventory Operations:")

	// Find existing product
	productResult := findProduct("P001")
	productResult.Map(func(product Product) Product {
		fmt.Printf("  Found product: %s ($%.2f, Stock: %d)\n",
			product.Name, product.Price, product.Stock)
		return product
	})

	// Try to find non-existent product
	missingResult := findProduct("P999")
	if missingResult.IsErr() {
		fmt.Printf("  Error: %v\n", missingResult.Error())
	}

	// Get active products
	activeProducts := inventory.Filter(func(id string, product Product) bool {
		return product.Active && product.Stock > 0
	})
	fmt.Printf("  Active products in stock: %d\n", activeProducts.Size())

	// Apply discount
	discountedInventory := inventory.Map(func(product Product) Product {
		return Product{
			ID:     product.ID,
			Name:   product.Name,
			Price:  product.Price * 0.85, // 15% discount
			Stock:  product.Stock,
			Active: product.Active,
		}
	})

	if product, found := discountedInventory.Get("P002"); found {
		fmt.Printf("  Discounted price for %s: $%.2f\n", product.Name, product.Price)
	}

	// Demonstrating immutability
	fmt.Println("\nImmutability Demonstration:")
	fmt.Printf("  Original inventory size: %d\n", inventory.Size())

	modifiedInventory := inventory.
		Set("P004", Product{
			ID:     "P004",
			Name:   "Smart Watch",
			Price:  299.99,
			Stock:  75,
			Active: true,
		}).
		Delete("P003")

	fmt.Printf("  Modified inventory size: %d\n", modifiedInventory.Size())
	fmt.Printf("  Original inventory unchanged: %d items\n", inventory.Size())
}

// main is the entry point for the examples program
func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "all":
		runAllExamples()
	case "core":
		runCoreExamples()
	case "enhanced":
		runEnhancedExamples()
	case "option":
		RunOptionExample()
	case "result":
		RunResultExample()
	case "iterator":
		RunIteratorExample()
	case "pattern":
		RunPatternExample()
	case "inventory":
		RunProductInventoryExample()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Println("=== RustGo Examples Runner ===")
	fmt.Println()
	fmt.Println("Usage: go run examples.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  all           - Run all examples")
	fmt.Println("  core          - Run core examples (Option, Result, Iterator, Pattern)")
	fmt.Println("  enhanced      - Run enhanced examples (Error handling, Immutable, Traits)")
	fmt.Println("  option        - Run Option type example")
	fmt.Println("  result        - Run Result type example")
	fmt.Println("  iterator      - Run Iterator example")
	fmt.Println("  pattern       - Run Pattern Matching example")
	fmt.Println("  inventory     - Run Product Inventory example (combined features)")
	fmt.Println("  help          - Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run examples.go all")
	fmt.Println("  go run examples.go option")
	fmt.Println("  go run examples.go enhanced")
	fmt.Println("  go run examples.go inventory")
}

func runAllExamples() {
	fmt.Println("Running all RustGo examples...")
	fmt.Println("==============================")
	fmt.Println()

	runCoreExamples()
	fmt.Println()
	runEnhancedExamples()
	fmt.Println()
	RunProductInventoryExample()
}

func runCoreExamples() {
	RunCoreExamples()
}

func runEnhancedExamples() {
	RunEnhancedExamples()
}
