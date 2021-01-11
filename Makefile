.PHONY: clean

REPO= github.com/ca-gip/dploy
NAME= dploy

dependency:
	go mod download

test: dependency
	GOARCH=amd64 go test -v ./internal/services

linux: test
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ./build/linux_amd64 -i $(GOPATH)/src/$(REPO)/main.go

darwin: test
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ./build/darwin_amd64 -i $(GOPATH)/src/$(REPO)/main.go

build: linux darwin