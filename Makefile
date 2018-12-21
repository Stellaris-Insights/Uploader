ifndef GOPATH
	$(error GOPATH not set. Please set GOPATH. \
	For more information about the GOPATH environment variable, see https://golang.org/doc/code.html#GOPATH)
endif

setup:
	go get -u github.com/golang/dep/cmd/dep
	wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.12.3

cisetup: setup
	go get github.com/mattn/goveralls

install:
	dep ensure

lint:
	golangci-lint run ./...

test:
	go test -v --race -cover ./api . ./cli/... -coverprofile tests.coverage

coverage:
	$(GOPATH)/bin/goveralls -service=travis-ci -coverprofile=tests.coverage

release:
	curl -sL http://git.io/goreleaser | bash