package structor

import (
	"fmt"
	"reflect"
)

// Set is a collection of codecs, for a specific tag name.
type Set struct {
	tagName string
	codecs  map[reflect.Type]*Codec
}

// NewSet returns a new codec set for the passed-in tag name.
func NewSet(tagName string) *Set {
	return &Set{tagName, make(map[reflect.Type]*Codec)}
}

// Add creates a new codec from the passed-in value and adds it to the set.
// It expects a struct, struct pointer or reflect.Type of a struct.
func (s *Set) Add(src interface{}) error {
	rType, err := codecType(src)
	if err != nil {
		return err
	}

	codec, err := newCodec(rType, s)
	if err != nil {
		return err
	}

	s.codecs[codec.rType] = codec
	return nil
}

// AddMust creates a new codec from the passed-in value and adds it to the set.
// It panics if an error occurs.
func (s *Set) AddMust(src interface{}) {
	err := s.Add(src)
	if err != nil {
		panic(err)
	}
}

// Get returns the according codec for the passed-in value.
// It returns an error if the value is invalid or no codec was found.
func (s *Set) Get(src interface{}) (*Codec, error) {
	rType, err := codecType(src)
	if err != nil {
		return nil, err
	}

	codec := s.codecs[rType]
	if codec == nil {
		return nil, fmt.Errorf("structor: no registered codec found for type '%v'", rType)
	}

	return codec, nil
}

// NewReader returns a new reader for the passed-in value, expecting a
// struct or struct pointer. It returns an error if the value is invalid
// or no matching codec was found.
func (s *Set) NewReader(src interface{}) (*Reader, error) {
	codec, err := s.Get(src)
	if err != nil {
		return nil, err
	}
	return newReader(src, codec), nil
}

// NewWriter returns a new writer for the passed-in value, expecting a
// struct pointer. It returns an error if the value is invalid or no
// matching codec was found.
func (s *Set) NewWriter(dst interface{}) (*Writer, error) {
	reader, err := s.NewReader(dst)
	if err != nil {
		return nil, err
	}
	return newWriter(dst, reader)
}

func codecType(src interface{}) (reflect.Type, error) {
	rVal := reflect.ValueOf(src)
	rKind := rVal.Kind()

	if t, ok := src.(reflect.Type); ok {
		return t, nil
	} else if rKind == reflect.Struct {
		return rVal.Type(), nil
	} else if rKind == reflect.Ptr && rVal.Elem().Kind() == reflect.Struct {
		return rVal.Elem().Type(), nil
	}

	return nil, fmt.Errorf("structor: value is not a struct, struct pointer or reflect.Type, but %q", rKind)
}
