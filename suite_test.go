package structor

import (
	"testing"
	. "github.com/101loops/bdd"
)

func TestSuite(t *testing.T) {
	RunSpecs(t, "reflector Suite")
}

func newTestSet() *Set {
	set := NewSet("test")
	set.AddMust(SimpleStruct{})
	return set
}

type SimpleStruct struct {
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:",omitempty"`
	_Mummy     uint64
	unexported uint64
}

func newSimpleStruct() SimpleStruct {
	return SimpleStruct{
		Dummy: "test",
		Yummy: 42,
	}
}

type ComplexStruct struct {
	SimpleStruct
	One   SimpleStruct
	Two   *SimpleStruct
	Three []SimpleStruct
	Four  map[string]SimpleStruct
}

func newComplexStruct() ComplexStruct {
	two := newSimpleStruct()
	return ComplexStruct{
		One: newSimpleStruct(),
		Two: &two,
	}
}

type RecursiveStruct struct {
	Level    int
	Parent   *RecursiveStruct
	Children []RecursiveStruct
}

func newRecursiveStruct() RecursiveStruct {
	return RecursiveStruct{
		Level:  0,
		Parent: &RecursiveStruct{Level: 1},
	}
}
