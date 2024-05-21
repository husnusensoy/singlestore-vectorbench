now := $(shell date +"%Y-%m-%d_%H:%M:%S")
sha1 := $(shell git rev-parse --short HEAD)
tag := $(shell git describe --tags --abbrev=0)
gov := $(shell go version)
hostname := $(shell hostname)

all: vectbench_linux_amd64-latest vectbench_darwin_amd64-latest

vectbench_darwin_amd64-latest: bin/vectbench_darwin_amd64-$(tag)-$(sha1)
vectbench_linux_amd64-latest: bin/vectbench_linux_amd64-$(tag)-$(sha1)

bin/vectbench_darwin_amd64-$(tag)-$(sha1):  src/vectbench/app.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.hostname=$(hostname) -X 'main.goV=$(gov)' -X main.tag=$(tag) -X main.sha1ver=$(sha1) -X main.buildTime=$(now)" -o $@  $<


bin/vectbench_linux_amd64-$(tag)-$(sha1):  src/vectbench/app.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.hostname=$(hostname) -X 'main.goV=$(gov)' -X main.tag=$(tag) -X main.sha1ver=$(sha1) -X main.buildTime=$(now)" -o $@  $<



