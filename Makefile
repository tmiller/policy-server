BINARY = policy-server
GOARCH = amd64

VERSION?=$(shell git describe HEAD)
BUILD=$(shell git rev-parse HEAD)
PACKAGE=${BINARY}_${VERSION}_${GOARCH}.deb

LDFLAGS = -ldflags "-s -w -X main.VERSION=${VERSION} -X main.BUILD=${BUILD}"

all: clean prep linux darwin package

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o build/linux-${GOARCH}/${BINARY} .

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o build/darwin-${GOARCH}/${BINARY} .

package:
	fpm \
		-f \
		-s dir \
		-t deb \
		-n ${BINARY} \
		-v ${VERSION} \
		-p build \
		--deb-systemd resources/policy-server.service \
		build/linux-${GOARCH}/${BINARY}=/opt/policy-server/${BINARY} \
		resources/crossdomain.xml=/opt/policy-server/crossdomain.xml \
		resources/env=/opt/policy-server/env \
		resources/policy-server-monitor.sh=/opt/policy-server/policy-server-monitor.sh

fmt:
	go fmt

prep:
	mkdir -p build/darwin-${GOARCH}
	mkdir -p build/linux-${GOARCH}

clean:
	- rm -r build

.PHONY: linux darwin fmt clean prep package
