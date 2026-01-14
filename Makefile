APP=yuan
VERSION=1.5.8
GITBRANCH=`git symbolic-ref --short -q HEAD`
GITREVISION=`git log -n1 --format=%H`
BUILDUSER=luyu151111@gamil.com
BUILDDATE=`date "+%Y-%m-%d %H:%M:%S"`
PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

## mod config
mod-config:
	@go env -w GO111MODULE='on'
	@go env -w GOPROXY='https://goproxy.cn,direct'

mod-tidy: mod-config
	@go mod tidy

mod-vendor: mod-config
	@go mod download
	@go mod vendor

## swag build
swag-build:
	@swag init -parseDependency --parseInternal  --parseDepth=1

## list: list packages and go files
list:
	@echo $(PACKAGES)
	@echo $(VETPACKAGES)
	@echo $(GOFILES)

## vet: static check
vet: mod-vendor
	@go vet $(VETPACKAGES)

## fmt: format code
fmt:
	@gofmt -s -w $(GOFILES)

## fmt-check: format result check
fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
	  echo "Please run 'make fmt' and commit the result:"; \
	  echo "$${diff}"; \
	  exit 1; \
	fi;

## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean
	@if [ -d build ] ; then rm -rf build ; fi

## bindata: package static resources
bindata:
	@go-bindata -o asset/asset.go -pkg=asset views/ docs/swagger.*

## test: runs go test with default values
test: bindata
	@go test ./... -v

## build: build the application to registry
build: clean test
	@go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)

## run: runs go run main.go
run:
	@go run -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' " main.go

## init-db: runs go run main.go -i
init-db:
	@LOGGER_OUT_STYLE=stdout go run -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' " main.go -i

## compile: build the application of different operating systems
compile:
	@# 32-Bit Systems
	@# FreeBDS
	@GOOS=freebsd GOARCH=386 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-freebsd-386 main.go
	@# Linux
	@GOOS=linux GOARCH=386 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-linux-386 main.go
	@# Windows
	@GOOS=windows GOARCH=386 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-windows-386 main.go
	@# 64-Bit
	@# FreeBDS
	@GOOS=freebsd GOARCH=amd64 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-freebsd-amd64 main.go
	@# MacOS
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-darwin-amd64 main.go

compile-asset: clean
	@# Linux 64-Bit
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-$(VERSION)-linux-amd64 main.go
	@# Windows 64-Bit
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o build/$(APP)-$(VERSION)-windows-amd64 main.go

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' Makefile | column -t -s ':' |  sed -e 's/^/ /'

## before push code: vet、fmt、fmt-check
before-push-code: list mod-tidy mod-vendor vet fmt fmt-check swag-build bindata

## all: execute test、build
all: vet fmt fmt-check test build

## init:
init: swag-build bindata mod-config mod-vendor

.PHONY: before-push-code all
