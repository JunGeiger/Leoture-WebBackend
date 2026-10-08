.PHONY: swagger
swagger:
	swag init -g leoture.go -o docs/swagger -d cmd,internal/handler,internal/errors,internal/types/request,internal/types/response,internal/model

.PHONY: dev
dev: swagger
	go run cmd/leoture.go

.PHONY: build
build: swagger
	go build -o bin/leotureweb cmd/leoture.go