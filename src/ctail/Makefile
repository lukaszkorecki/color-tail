GOPATH=`pwd`/../..
main = ctail.go
# find all package names in src and add them to list
test_packages=`find ./src/ -type d | sed 's/^.*src.*\///'`

build:
	@echo Building in $(GOPATH)
	@env GOPATH=$(GOPATH) go build -v $(main)

test:
	@env GOPATH=$(GOPATH) go test -v $(test_packages)

dependencies:
	@env GOPATH=$(GOPATH) go get github.com/howeyc/fsnotify
