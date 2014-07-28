package structor

import (
	"testing"
	. "github.com/101loops/bdd"
)

type SimpleStruct struct {
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:",omitempty"`
	Ignored    uint64 `test:"-"`
	unexported uint64
}

type ComplexStruct struct {
	One   SimpleStruct
	Two   *SimpleStruct
	Three []SimpleStruct
	Four  map[string]SimpleStruct
}

type RecursiveStruct struct {
	Level    int
	Parent   *RecursiveStruct
	Children []RecursiveStruct
}

func TestSuite(t *testing.T) {
	RunSpecs(t, "reflector Suite")
}

func newTestSet() *Set {
	set := NewSet("test")
	set.AddMust(SimpleStruct{})
	return set
}

func newSimpleStruct() SimpleStruct {
	return SimpleStruct{
		Dummy: "test",
		Yummy: 42,
	}
}

func newComplexStruct() ComplexStruct {
	two := newSimpleStruct()
	return ComplexStruct{
		One: newSimpleStruct(),
		Two: &two,
	}
}

func newRecursiveStruct() RecursiveStruct {
	return RecursiveStruct{
		Level:  0,
		Parent: &RecursiveStruct{Level: 1},
	}
}
