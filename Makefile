.PHONY: default di clean build

default: di build;

di:
	go run github.com/google/wire/cmd/wire ./cmd/enhance_api

clean:
	rm -rf ./out

build:
	go build -o out/enhance_api ./cmd/enhance_api

install:
	go install ./cmd/enhance_api

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./out/enhance_api ./cmd/enhance_api