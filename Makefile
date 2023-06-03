build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test -v -count=1 ./...
start: build
	docker-compose up --build sensor-control

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

lint:
	echo "Starting linters"
	golangci-lint run ./...

PHONY: generate
generate:
	mkdir -p pkg/models/
	protoc --go_out=./pkg/models --go_opt=paths=source_relative \
		   api/models/snapshot.proto  
	mv pkg/models/api/models/* pkg/models
	rm -rf pkg/models/api