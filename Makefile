.PHONY: gen
gen:
	protoc --proto_path=internal/controller/grpcapi/proto/blacklist internal/controller/grpcapi/proto/blacklist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/whitelist internal/controller/grpcapi/proto/whitelist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/bucket internal/controller/grpcapi/proto/bucket/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/authorization internal/controller/grpcapi/proto/authorization/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import

.PHONY: mock.gen
mock.gen:
	mockgen -source=internal/domain/service/blacklist.go -destination=internal/store/postgressql/adapters/mocks/mock_blacklist.go
	mockgen -source=internal/domain/service/whitelist.go -destination=internal/store/postgressql/adapters/mocks/mock_whitelist.go

.PHONY: clean
clean:
	rm -f internal/controller/grpcapi/blacklistpb/*
	rm -f internal/controller/grpcapi/whitelistpb/*
	rm -f internal/controller/grpcapi/bucketpb/*
	rm -f internal/controller/grpcapi/authorization/*

.PHONY: build.docker
build.docker:
	docker build --tag abf --  .

.PHONY: run.docker
run.docker:
	docker run -p 8080:8080 -it --name abf_container abf

.PHONY: build
build:
	docker-compose build

.PHONY: build.bin
build.bin:
	go build -o ./build/anti_bruteforce_app/anti_bruteforce ./cmd/server

.PHONY: run
run: build
	docker-compose up

.PHONY: stop
stop:
	docker-compose down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm 	coverage.out

.PHONY: test
test:
	go test -race ./...

format-go:
	golangci-lint cache clean
	golangci-lint run --fix ./...