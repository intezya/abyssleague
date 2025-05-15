ent generate command
go run -mod=mod entgo.io/ent/cmd/ent generate --target=./internal/infrastructure/ent ./internal/infrastructure/ent/schema

swagger documentation (it also must be executed in CI)
swag init --dir ./cmd,./internal/adapters/controller/http/handlers,./docs/examples --parseDependency --output ./docs
swag fmt (should be used as pre-commit hook)

c:\bin\golangci-lint run ./... --default=all --fix
