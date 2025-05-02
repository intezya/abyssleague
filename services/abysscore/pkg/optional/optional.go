package optional

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

func NewP[T any](value *T) Optional[T] {
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

func (e Optional[T]) IsSet() bool {
	return e.isSet && e.value != nil
}

func (e Optional[T]) Value() (value T, ok bool) {
	if !e.isSet || e.value == nil {
		var zero T

		return zero, false
	}

	return *e.value, true
}

func (e Optional[T]) MustValue() T {
	if !e.isSet || e.value == nil {
		panic("attempted to get value from empty or nil optional")
	}

	return *e.value
}

func (e Optional[T]) Default(defaultValue T) T {
	if e.isSet && e.value != nil {
		return *e.value
	}

	return defaultValue
}

// DefaultFn provides method for getting default value from fn without extra call if isSet.
func (e Optional[T]) DefaultFn(defaultValue func() T) T {
	if e.isSet && e.value != nil {
		return *e.value
	}

	return defaultValue()
}
