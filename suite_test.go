package structor

import (
	. "github.com/101loops/bdd"
	"testing"
)

type TestStruct struct {
	Dummy      string `test:"dummytag"`
	Yummy      int    `test:",omitempty"`
	Ignored    uint64 `test:"-"`
	unexported uint64
}

func TestSuite(t *testing.T) {
	RunSpecs(t, "reflector Suite")
}

func newTestSet() *Set {
	set := NewSet("test")
	set.AddMust(TestStruct{})
	return set
}

func newTestData() TestStruct {
	return TestStruct{
		Dummy: "test",
		Yummy: 42,
	}
}
