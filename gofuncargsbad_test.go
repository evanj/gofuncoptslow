package main

import (
	"os"
	"testing"
)

// Foo is an example type with a constructor that takes options. It counts the number
// of options that were supplied to the constructor.
type Foo struct {
	optionCount int
}

type FooOption func(f *Foo)

func WithBoolOption() FooOption {
	return func(f *Foo) {
		f.optionCount++
	}
}

func WithIntOption(x int) FooOption {
	return func(f *Foo) {
		if x > 0 {
			f.optionCount++
		}
	}
}

func NewFoo(options ...FooOption) *Foo {
	f := &Foo{}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type fooFactory interface {
	newFoo(options ...FooOption) *Foo
	newFooStruct(config FooConfig) *Foo
}

type wrapperFooFactory struct{}

func (g *wrapperFooFactory) newFoo(options ...FooOption) *Foo {
	return NewFoo(options...)
}

func (g *wrapperFooFactory) newFooStruct(config FooConfig) *Foo {
	return NewFooStruct(config)
}

// newFactory creates a fooFactory based on an environment variable, so the underlying type cannot
// be statically determined. This means the compiler has to call the functions via the interface.
func newFactory() fooFactory {
	if len(os.Getenv("BUG")) > 0 {
		return nil
	}
	return &wrapperFooFactory{}
}

type FooConfig struct {
	BoolOption bool
	IntOption  int
}

func NewFooStruct(config FooConfig) *Foo {
	f := &Foo{}
	if config.BoolOption {
		f.optionCount++
	}
	if config.IntOption > 0 {
		f.optionCount++
	}
	return f
}

// forbid inlining for an "apples to apples" performance comparison
//go:noinline
func NewFooStructNoInline(config FooConfig) *Foo {
	return NewFooStruct(config)
}

func CallNewFoo() *Foo {
	return NewFoo(WithBoolOption(), WithIntOption(42))
}

func CallNewFooStruct() *Foo {
	return NewFooStruct(FooConfig{BoolOption: true, IntOption: 42})
}

func CallNewFooStructNoInline() *Foo {
	return NewFooStructNoInline(FooConfig{BoolOption: true, IntOption: 42})
}

func CallNewFooInterface(factory fooFactory) *Foo {
	return factory.newFoo(WithBoolOption(), WithIntOption(42))
}

func CallNewFooStructInterface(factory fooFactory) *Foo {
	return factory.newFooStruct(FooConfig{BoolOption: true, IntOption: 42})
}

func BenchmarkNewFoo(b *testing.B) {
	// these calls are inlined by the compiler
	// 1 allocs/op: foo is allocated
	b.Run("no_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFoo()
			if f.optionCount != 0 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("1_option_no_arguments", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFoo(WithBoolOption())
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("1_option_with_argument", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFoo(WithIntOption(i + 1))
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("2_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFoo(WithBoolOption(), WithIntOption(i+1))
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	// These versions all get inlined and effectively do nothing, with ZERO allocations
	// even foo does not get allocated
	b.Run("struct_no_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFooStruct(FooConfig{})
			if f.optionCount != 0 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_1_option_no_arguments", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFooStruct(FooConfig{BoolOption: true})
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_1_option_with_argument", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFooStruct(FooConfig{IntOption: i + 1})
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_2_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := NewFooStruct(FooConfig{BoolOption: true, IntOption: i + 1})
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	// create factory with a non-deterministic type
	// Due to the interface type, these functions cannot be inlined:
	// needs to alloc the ... argument, and maybe the function closures
	factory := newFactory()
	b.Run("interface_no_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFoo()
			if f.optionCount != 0 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("interface_1_option_no_arguments", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFoo(WithBoolOption())
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("interface_1_option_with_argument", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFoo(WithIntOption(i + 1))
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("interface_2_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFoo(WithBoolOption(), WithIntOption(i+1))
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	// same factory type with a struct value argument:
	b.Run("struct_interface_no_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFooStruct(FooConfig{})
			if f.optionCount != 0 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_interface_1_option_no_arguments", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFooStruct(FooConfig{BoolOption: true})
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_interface_1_option_with_argument", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFooStruct(FooConfig{IntOption: i + 1})
			if f.optionCount != 1 {
				b.Fatal("BUG")
			}
		}
	})
	b.Run("struct_interface_2_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := factory.newFooStruct(FooConfig{BoolOption: true, IntOption: i + 1})
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})
}

func BenchmarkCallNewFoo(b *testing.B) {
	// 1 allocs/op: foo is allocated
	b.Run("func_options", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := CallNewFoo()
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	// 0 allocs/op: foo is not allocated
	b.Run("struct", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := CallNewFooStruct()
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	// 1 allocs/op: foo is allocated
	b.Run("struct_noinline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := CallNewFooStructNoInline()
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	factory := newFactory()
	b.Run("func_options_interface", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := CallNewFooInterface(factory)
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})

	b.Run("struct_interface", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := CallNewFooStructInterface(factory)
			if f.optionCount != 2 {
				b.Fatal("BUG")
			}
		}
	})
}
