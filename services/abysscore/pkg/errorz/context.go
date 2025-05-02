package errorz

type Context interface {
	Path() string
	Status(code int) Context
	JSON(data interface{}) error
}
