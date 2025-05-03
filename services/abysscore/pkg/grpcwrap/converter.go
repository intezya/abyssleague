package grpcwrap

import (
	"fmt"
)

type SimpleConverter[From, To any] struct {
	ConvertFunc func(from From) (To, error)
}

func (c *SimpleConverter[From, To]) Convert(from From) (To, error) {
	return c.ConvertFunc(from)
}

type SliceConverter[From, To any] struct {
	ElementConverter TypeConverter[From, To]
}

func (c *SliceConverter[From, To]) Convert(fromSlice []From) ([]To, error) {
	result := make([]To, 0, len(fromSlice))

	for _, item := range fromSlice {
		converted, err := c.ElementConverter.Convert(item)
		if err != nil {
			return nil, fmt.Errorf("failed to convert slice element: %w", err)
		}

		result = append(result, converted)
	}

	return result, nil
}
