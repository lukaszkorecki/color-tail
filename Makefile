GOPATH=`pwd`
VERSION := "v`cat VERSION`"
main = ctail.go
# find all package names in src and add them to list
test_packages=`find -type d | egrep -v "src|.git|.pkg"`

all: dependencies test build

build:
	@echo Building in $(GOPATH)
	@env GOPATH=$(GOPATH) go build -v $(main)

test:
	@echo Testing!
	@env GOPATH=$(GOPATH) go test -v $(test_packages)

dependencies:
	@echo installing dependencies
	@mkdir -p src
	@env GOPATH=$(GOPATH) go get github.com/howeyc/fsnotify

edit:
	@env GOPATH=$(GOPATH) vim .

release:
	git tag $(VERSION)
	git commit -m "Release: $(VERSION)
	@echo $(VERSION) is ready to push
