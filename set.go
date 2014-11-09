package structor

import (
	"fmt"
	"reflect"
)

// Set is a collection of codecs, for a specific tag name.
type Set struct {
	tagName  string
	codecs   map[reflect.Type]*Codec
	validate func(*Set, *Codec) error
}

// NewSet returns a new codec set for the passed-in tag name.
func NewSet(tagName string) *Set {
	return &Set{
		tagName: tagName,
		codecs:  make(map[reflect.Type]*Codec),
		validate: func(*Set, *Codec) error {
			return nil
		},
	}
}

// Add creates a new codec from the passed-in value and adds it to the set.
// It expects a struct, struct pointer or reflect.Type of a struct.
func (s *Set) Add(src interface{}) error {
	rType, err := typeOf(src)
	if err != nil {
		return err
	}

	if _, found := s.codecs[rType]; found {
		return nil
	}

	codec := newCodec(rType, s.tagName)
	s.codecs[rType] = codec // added eagerly for recursive types

	for _, f := range codec.Fields() {
		subTypes := []*reflect.Type{&f.Type, f.KeyType, f.ElemType}
		for _, typ := range subTypes {
			if typ == nil {
				continue
			}

			var subType reflect.Type
			switch (*typ).Kind() {
			case reflect.Struct:
				subType = *typ
			case reflect.Ptr:
				if (*typ).Elem().Kind() == reflect.Struct {
					subType = (*typ).Elem()
				}
			}
			if subType == nil {
				continue
			}

			if err := s.Add(subType); err != nil {
				delete(s.codecs, rType)
				return err
			}
		}
	}

	if err := s.validate(s, codec); err != nil {
		delete(s.codecs, rType)
		return err
	}

	codec.complete = true
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

// SetValidateFunc defines a function to validate a codec before it is added.
// When it returns an error the codec is not added to the set.
func (s *Set) SetValidateFunc(fn func(*Set, *Codec) error) {
	s.validate = fn
}

// Get returns the according codec for the passed-in value.
// It returns an error if the value is invalid or no codec was found.
func (s *Set) Get(src interface{}) (*Codec, error) {
	rType, err := typeOf(src)
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

func typeOf(src interface{}) (reflect.Type, error) {
	rVal := reflect.ValueOf(src)
	rKind := rVal.Kind()

	if t, ok := src.(reflect.Type); ok {
		return t, nil
	} else if rKind == reflect.Struct {
		return rVal.Type(), nil
	} else if rKind == reflect.Ptr && rVal.Elem().Kind() == reflect.Struct {
		return rVal.Elem().Type(), nil
	}

	return nil, fmt.Errorf("structor: value is not a struct, struct pointer or reflect.Type - but %q", rKind)
}
