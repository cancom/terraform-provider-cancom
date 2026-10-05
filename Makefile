TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOSTNAME=cancom.de
NAMESPACE=cancom
NAME=cancom
BINARY=terraform-provider-${NAME}
VERSION=0.0.3
OS_ARCH=darwin_arm64

default: install

BIN=$(CURDIR)/bin
$(BIN)/%:
	@echo "Installing tools from tools/tools.go"
	@cat tools/tools.go | grep _ | awk -F '"' '{print $$2}' | GOBIN=$(BIN) xargs -tI {} go install {}

build: clean
	CGO_ENABLED=0 go build -o $(BIN)/${BINARY}
	@sh -c "'$(CURDIR)/hack/generate-dev-overrides.sh'"

release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp $(BIN)/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	CGO_ENABLED=0 go test $(TEST) -timeout=30s -parallel=4

testacc: 
	CGO_ENABLED=0 TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	gofmt -w $(GOFMT_FILES)

generate-docs: $(BIN)/tfplugindocs
	$(BIN)/tfplugindocs generate

validate-docs: $(BIN)/tfplugindocs
	$(BIN)/tfplugindocs validate

clean:
	rm -rf ./bin

.PHONY: build testacc fmt validate-docs generate-docs clean
