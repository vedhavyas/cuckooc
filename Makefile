install:
	go get -u github.com/Masterminds/glide
	glide install

test:
	go test $(shell go list ./... | grep -v '/vendor/') --cover

package:
	go clean
	OS="darwin"
	CGO_ENABLED=0 GOOS=$$OS go build ./cmd/cuckooc/

all: install test package
