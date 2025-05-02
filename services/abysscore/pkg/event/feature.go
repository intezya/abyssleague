package eventlib

import (
	"fmt"
	"runtime/debug"
	"time"
)

func RecoverMiddleware(log func(string, ...interface{})) Middleware {
	return func(next Handler) Handler {
		return func(event ApplicationEvent) {
			defer func() {
				if rec := recover(); rec != nil {
					stackTrace := debug.Stack()
					timestamp := time.Now()

					log(
						"panic recovered",
						"error", fmt.Sprintf("%v", rec),
						"stacktrace", string(stackTrace),
						"timestamp", timestamp,
					)
				}
			}()
			next(event)
		}
	}
}

func LoggingMiddleware(next Handler, log func(args ...interface{})) Handler {
	return func(event ApplicationEvent) {
		log("handling event: ", typeName(event))
		next(event)
	}
}

func LoggingfMiddleware(next Handler, log func(template string, args ...interface{})) Handler {
	return func(event ApplicationEvent) {
		log("handling event: %s", typeName(event))
		next(event)
	}
}
