BINARY = policy-server
GOARCH = amd64

all: clean prep linux darwin

linux: 
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o build/linux-${GOARCH}/${BINARY} .

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o build/darwin-${GOARCH}/${BINARY} .

fmt:
	go fmt

prep:
	mkdir -p build/darwin-${GOARCH}
	mkdir -p build/linux-${GOARCH}

clean:
	- rm -r build

.PHONY: linux darwin fmt clean prep
