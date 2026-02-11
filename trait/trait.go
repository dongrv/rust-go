// Package trait provides Rust-like trait system for Go with compile-time polymorphism
// and better code organization through interface composition.
package trait

import (
	"fmt"
	"reflect"
)

// Trait is a marker interface for all traits
type Trait interface {
	// traitName returns the name of the trait
	traitName() string
}

// TraitRegistry maintains a registry of trait implementations
type TraitRegistry struct {
	implementations map[string]map[reflect.Type]interface{}
}

var globalRegistry = &TraitRegistry{
	implementations: make(map[string]map[reflect.Type]interface{}),
}

// Register registers a trait implementation for a specific type
func Register[T Trait, Impl any](trait T, implementation Impl) {
	traitName := trait.traitName()
	typeKey := reflect.TypeOf((*Impl)(nil)).Elem()

	if globalRegistry.implementations[traitName] == nil {
		globalRegistry.implementations[traitName] = make(map[reflect.Type]interface{})
	}

	globalRegistry.implementations[traitName][typeKey] = implementation
}

// Get retrieves a trait implementation for a specific type
func Get[T Trait, Impl any](trait T) (Impl, bool) {
	traitName := trait.traitName()
	typeKey := reflect.TypeOf((*Impl)(nil)).Elem()

	if impls, ok := globalRegistry.implementations[traitName]; ok {
		if impl, ok := impls[typeKey]; ok {
			return impl.(Impl), true
		}
	}

	var zero Impl
	return zero, false
}

// Implementor represents a type that implements one or more traits
type Implementor struct {
	value      interface{}
	traitImpls map[string]interface{}
}

// NewImplementor creates a new Implementor for the given value
func NewImplementor(value interface{}) *Implementor {
	return &Implementor{
		value:      value,
		traitImpls: make(map[string]interface{}),
	}
}

// With adds a trait implementation to the implementor
func (i *Implementor) With(traitName string, implementation interface{}) *Implementor {
	i.traitImpls[traitName] = implementation
	return i
}

// GetTrait retrieves a trait implementation from the implementor
func (i *Implementor) GetTrait(traitName string) (interface{}, bool) {
	impl, ok := i.traitImpls[traitName]
	return impl, ok
}

// Value returns the underlying value
func (i *Implementor) Value() interface{} {
	return i.value
}

// TraitObject represents a type-erased trait object (dynamic dispatch)
type TraitObject struct {
	data   interface{}
	vtable map[string]interface{}
}

// NewTraitObject creates a new trait object
func NewTraitObject(data interface{}, vtable map[string]interface{}) *TraitObject {
	return &TraitObject{
		data:   data,
		vtable: vtable,
	}
}

// Call calls a method on the trait object
func (to *TraitObject) Call(methodName string, args ...interface{}) ([]interface{}, error) {
	method, ok := to.vtable[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found in vtable", methodName)
	}

	methodValue := reflect.ValueOf(method)
	if methodValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("vtable entry for %s is not a function", methodName)
	}

	// Prepare arguments
	in := make([]reflect.Value, len(args)+1)
	in[0] = reflect.ValueOf(to.data)
	for i, arg := range args {
		in[i+1] = reflect.ValueOf(arg)
	}

	// Call the method
	results := methodValue.Call(in)

	// Convert results to interface{}
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out, nil
}

// Display is a trait for types that can be displayed as strings
type Display interface {
	Trait
	Display() string
}

// displayTrait is the concrete trait type
type displayTrait struct{}

func (d displayTrait) traitName() string {
	return "Display"
}

func (d displayTrait) Display() string {
	return "Display trait"
}

// DisplayTrait is the singleton Display trait
var DisplayTrait Display = displayTrait{}

// Debug is a trait for types that can provide debug information
type Debug interface {
	Trait
	Debug() string
}

// debugTrait is the concrete trait type
type debugTrait struct{}

func (d debugTrait) traitName() string {
	return "Debug"
}

func (d debugTrait) Debug() string {
	return "Debug trait"
}

// DebugTrait is the singleton Debug trait
var DebugTrait Debug = debugTrait{}

// Clone is a trait for types that can be cloned
type Clone interface {
	Trait
	Clone() interface{}
}

// cloneTrait is the concrete trait type
type cloneTrait struct{}

func (c cloneTrait) traitName() string {
	return "Clone"
}

func (c cloneTrait) Clone() interface{} {
	return cloneTrait{}
}

// CloneTrait is the singleton Clone trait
var CloneTrait Clone = cloneTrait{}

// Eq is a trait for types that support equality comparison
type Eq interface {
	Trait
	Eq(other interface{}) bool
}

// eqTrait is the concrete trait type
type eqTrait struct{}

func (e eqTrait) traitName() string {
	return "Eq"
}

func (e eqTrait) Eq(other interface{}) bool {
	return true
}

// EqTrait is the singleton Eq trait
var EqTrait Eq = eqTrait{}

// Ord is a trait for types that support ordering
type Ord interface {
	Trait
	Cmp(other interface{}) int // -1 for less, 0 for equal, 1 for greater
}

// ordTrait is the concrete trait type
type ordTrait struct{}

func (o ordTrait) traitName() string {
	return "Ord"
}

func (o ordTrait) Cmp(other interface{}) int {
	return 0
}

// OrdTrait is the singleton Ord trait
var OrdTrait Ord = ordTrait{}

// Hash is a trait for types that can be hashed
type Hash interface {
	Trait
	Hash() uint64
}

// hashTrait is the concrete trait type
type hashTrait struct{}

func (h hashTrait) traitName() string {
	return "Hash"
}

func (h hashTrait) Hash() uint64 {
	return 0
}

// HashTrait is the singleton Hash trait
var HashTrait Hash = hashTrait{}

// Default is a trait for types that have a default value
type Default interface {
	Trait
	Default() interface{}
}

// defaultTrait is the concrete trait type
type defaultTrait struct{}

func (d defaultTrait) traitName() string {
	return "Default"
}

func (d defaultTrait) Default() interface{} {
	return defaultTrait{}
}

// DefaultTrait is the singleton Default trait
var DefaultTrait Default = defaultTrait{}

// Iterator is a trait for types that can be iterated over
type Iterator interface {
	Trait
	Next() (interface{}, bool)
}

// iteratorTrait is the concrete trait type
type iteratorTrait struct{}

func (i iteratorTrait) traitName() string {
	return "Iterator"
}

func (i iteratorTrait) Next() (interface{}, bool) {
	return nil, false
}

// IteratorTrait is the singleton Iterator trait
var IteratorTrait Iterator = iteratorTrait{}

// FromStr is a trait for types that can be parsed from strings
type FromStr interface {
	Trait
	FromStr(s string) (interface{}, error)
}

// fromStrTrait is the concrete trait type
type fromStrTrait struct{}

func (f fromStrTrait) traitName() string {
	return "FromStr"
}

func (f fromStrTrait) FromStr(s string) (interface{}, error) {
	return s, nil
}

// FromStrTrait is the singleton FromStr trait
var FromStrTrait FromStr = fromStrTrait{}

// ToString is a trait for types that can be converted to strings
type ToString interface {
	Trait
	ToString() string
}

// toStringTrait is the concrete trait type
type toStringTrait struct{}

func (t toStringTrait) traitName() string {
	return "ToString"
}

func (t toStringTrait) ToString() string {
	return "ToString trait"
}

// ToStringTrait is the singleton ToString trait
var ToStringTrait ToString = toStringTrait{}

// Derive is a helper for deriving traits automatically
type Derive struct {
	target interface{}
}

// NewDerive creates a new Derive helper for the target type
func NewDerive(target interface{}) *Derive {
	return &Derive{target: target}
}

// Display derives the Display trait
func (d *Derive) Display() *Derive {
	// Auto-derive Display using reflection
	targetType := reflect.TypeOf(d.target)
	impl := struct {
		DisplayFunc func() string
	}{
		DisplayFunc: func() string {
			return fmt.Sprintf("%v", d.target)
		},
	}
	// Register with the target type as key
	if globalRegistry.implementations["Display"] == nil {
		globalRegistry.implementations["Display"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Display"][targetType] = impl
	return d
}

// Debug derives the Debug trait
func (d *Derive) Debug() *Derive {
	// Auto-derive Debug using reflection
	targetType := reflect.TypeOf(d.target)
	impl := struct {
		DebugFunc func() string
	}{
		DebugFunc: func() string {
			return fmt.Sprintf("%#v", d.target)
		},
	}
	// Register with the target type as key
	if globalRegistry.implementations["Debug"] == nil {
		globalRegistry.implementations["Debug"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Debug"][targetType] = impl
	return d
}

// Clone derives the Clone trait
func (d *Derive) Clone() *Derive {
	// Auto-derive Clone using reflection
	targetType := reflect.TypeOf(d.target)
	impl := struct {
		CloneFunc func() interface{}
	}{
		CloneFunc: func() interface{} {
			val := reflect.ValueOf(d.target)
			if val.Kind() == reflect.Ptr {
				// For pointers, create a new pointer to a copy of the value
				elem := reflect.New(val.Elem().Type())
				elem.Elem().Set(val.Elem())
				return elem.Interface()
			}
			// For values, return a copy
			return reflect.New(val.Type()).Elem().Interface()
		},
	}
	// Register with the target type as key
	if globalRegistry.implementations["Clone"] == nil {
		globalRegistry.implementations["Clone"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Clone"][targetType] = impl
	return d
}

// Eq derives the Eq trait
func (d *Derive) Eq() *Derive {
	// Auto-derive Eq using reflection
	targetType := reflect.TypeOf(d.target)
	impl := struct {
		EqFunc func(other interface{}) bool
	}{
		EqFunc: func(other interface{}) bool {
			return reflect.DeepEqual(d.target, other)
		},
	}
	// Register with the target type as key
	if globalRegistry.implementations["Eq"] == nil {
		globalRegistry.implementations["Eq"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Eq"][targetType] = impl
	return d
}

// Default derives the Default trait
func (d *Derive) Default() *Derive {
	// Auto-derive Default using reflection
	targetType := reflect.TypeOf(d.target)
	impl := struct {
		DefaultFunc func() interface{}
	}{
		DefaultFunc: func() interface{} {
			t := reflect.TypeOf(d.target)
			return reflect.New(t).Elem().Interface()
		},
	}
	// Register with the target type as key
	if globalRegistry.implementations["Default"] == nil {
		globalRegistry.implementations["Default"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Default"][targetType] = impl
	return d
}

// TraitComposition allows composing multiple traits
type TraitComposition struct {
	traits []string
}

// Compose creates a new trait composition
func Compose(traits ...string) *TraitComposition {
	return &TraitComposition{traits: traits}
}

// Implement creates an implementor with all composed traits
func (tc *TraitComposition) Implement(value interface{}) *Implementor {
	impl := NewImplementor(value)
	for _, trait := range tc.traits {
		// Look up trait implementation in registry
		if impls, ok := globalRegistry.implementations[trait]; ok {
			for typeKey, traitImpl := range impls {
				if reflect.TypeOf(value).AssignableTo(typeKey) {
					impl.With(trait, traitImpl)
					break
				}
			}
		}
	}
	return impl
}

// TraitBound represents a trait bound for generic constraints
type TraitBound struct {
	traitName string
}

// NewBound creates a new trait bound
func NewBound(traitName string) *TraitBound {
	return &TraitBound{traitName: traitName}
}

// Check checks if a value satisfies the trait bound
func (tb *TraitBound) Check(value interface{}) bool {
	if impls, ok := globalRegistry.implementations[tb.traitName]; ok {
		valueType := reflect.TypeOf(value)
		for typeKey := range impls {
			if valueType.AssignableTo(typeKey) {
				return true
			}
		}
	}
	return false
}

// Require panics if the value doesn't satisfy the trait bound
func (tb *TraitBound) Require(value interface{}) {
	if !tb.Check(value) {
		panic(fmt.Sprintf("value of type %T does not satisfy trait bound %s", value, tb.traitName))
	}
}

// DynamicDispatch provides runtime polymorphism through trait objects
type DynamicDispatch struct {
	objects map[string]*TraitObject
}

// NewDynamicDispatch creates a new dynamic dispatch container
func NewDynamicDispatch() *DynamicDispatch {
	return &DynamicDispatch{
		objects: make(map[string]*TraitObject),
	}
}

// Add adds a trait object with the given name
func (dd *DynamicDispatch) Add(name string, obj *TraitObject) {
	dd.objects[name] = obj
}

// Call calls a method on a named trait object
func (dd *DynamicDispatch) Call(name, method string, args ...interface{}) ([]interface{}, error) {
	obj, ok := dd.objects[name]
	if !ok {
		return nil, fmt.Errorf("trait object %s not found", name)
	}
	return obj.Call(method, args...)
}

// TraitAlias creates an alias for a trait
func TraitAlias(original, alias string) {
	if impls, ok := globalRegistry.implementations[original]; ok {
		globalRegistry.implementations[alias] = impls
	}
}

// HasTrait checks if a type has a specific trait implementation
func HasTrait(traitName string, value interface{}) bool {
	if impls, ok := globalRegistry.implementations[traitName]; ok {
		valueType := reflect.TypeOf(value)
		for typeKey := range impls {
			if valueType.AssignableTo(typeKey) {
				return true
			}
		}
	}
	return false
}

// GetTraitNames returns all registered trait names
func GetTraitNames() []string {
	names := make([]string, 0, len(globalRegistry.implementations))
	for name := range globalRegistry.implementations {
		names = append(names, name)
	}
	return names
}

// ClearRegistry clears the trait registry (mainly for testing)
func ClearRegistry() {
	globalRegistry.implementations = make(map[string]map[reflect.Type]interface{})
}

// Example implementations for common types

func init() {
	// Register Display for int
	intType := reflect.TypeOf(0)
	if globalRegistry.implementations["Display"] == nil {
		globalRegistry.implementations["Display"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Display"][intType] = struct {
		DisplayFunc func() string
	}{
		DisplayFunc: func() string {
			return "int"
		},
	}

	// Register Display for string
	stringType := reflect.TypeOf("")
	globalRegistry.implementations["Display"][stringType] = struct {
		DisplayFunc func() string
	}{
		DisplayFunc: func() string {
			return "string"
		},
	}

	// Register Eq for int
	if globalRegistry.implementations["Eq"] == nil {
		globalRegistry.implementations["Eq"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Eq"][intType] = struct {
		EqFunc func(other interface{}) bool
	}{
		EqFunc: func(other interface{}) bool {
			if _, ok := other.(int); ok {
				return true // Simplified for example
			}
			return false
		},
	}

	// Register Clone for int
	if globalRegistry.implementations["Clone"] == nil {
		globalRegistry.implementations["Clone"] = make(map[reflect.Type]interface{})
	}
	globalRegistry.implementations["Clone"][intType] = struct {
		CloneFunc func() interface{}
	}{
		CloneFunc: func() interface{} {
			return 0
		},
	}

	// Debug: Print registered trait names
	// fmt.Println("Registered traits in init():")
	// for traitName := range globalRegistry.implementations {
	//     fmt.Printf("  - %s\n", traitName)
	// }
}
