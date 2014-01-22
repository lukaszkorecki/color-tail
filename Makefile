VERSION := "v"`head -1 VERSION`
main = main.go
# find all package names in src and add them to list
local_packages :=`find -type d | egrep -v "src|.git|.pkg"`

all: dependencies test build

build:
	@echo Building in $(GOPATH)
	go build -ldflags "-X main.version $(VERSION)"

test:
	@echo Testing!
	@env GOPATH=$(GOPATH) go test -v $(local_packages)

dependencies:
	@echo installing dependencies
	@mkdir -p src
	@env GOPATH=$(GOPATH) go get github.com/howeyc/fsnotify

fmt:
	@go fmt .
	@for d in $(local_packages) ; do echo $$d ; go fmt $$d/*.go ; done

release:
	git tag $(VERSION)
	git commit -m "Release: $(VERSION)
	@echo $(VERSION) is ready to push
