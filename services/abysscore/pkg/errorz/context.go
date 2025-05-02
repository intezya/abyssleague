package errorz

type Context interface {
	Path() string
	Status(int) Context
	JSON(data interface{}) error
}
