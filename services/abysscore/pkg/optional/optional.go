package optional

import jsoniter "github.com/json-iterator/go"

/*
	Optional package provides custom (extended) solution for optional types
*/

type Optional[T any] struct {
	isSet bool
	value *T
}

func New[T any](value T) Optional[T] {
	return Optional[T]{
		true,
		&value,
	}
}

func FromPtr[T any](value *T) Optional[T] {
	if value == nil {
		return Optional[T]{
			false,
			nil,
		}
	}

	return Optional[T]{
		true,
		value,
	}
}

// EmptyOptional returns a new Optional that does not have a value set.
func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{
		false,
		nil,
	}
}

func (o Optional[T]) IsSet() bool {
	return o.isSet && o.value != nil
}

func (o Optional[T]) Value() (value T, ok bool) {
	if !o.isSet || o.value == nil {
		var zero T

		return zero, false
	}

	return *o.value, true
}

func (o Optional[T]) ValueOrNil() *T {
	if !o.isSet {
		return nil
	}

	return o.value
}

func (o Optional[T]) MustValue() T {
	if !o.isSet || o.value == nil {
		panic("attempted to get value from empty or nil optional")
	}

	return *o.value
}

func (o Optional[T]) Default(defaultValue T) T {
	if o.isSet && o.value != nil {
		return *o.value
	}

	return defaultValue
}

// DefaultFn provides method for getting default value from fn without extra call if isSet.
func (o Optional[T]) DefaultFn(defaultValue func() T) T {
	if o.isSet && o.value != nil {
		return *o.value
	}

	return defaultValue()
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(o.value)
}
