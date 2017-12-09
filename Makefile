BINARY = policy-server
GOARCH = amd64

VERSION?=$(shell git describe HEAD)
BUILD=$(shell git rev-parse HEAD)
PACKAGE=${BINARY}-go_${VERSION}_${GOARCH}.deb

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
		-n ${BINARY}-go \
		-v ${VERSION} \
		-p build \
		--deb-upstart resources/policy-server-go.conf \
		--after-upgrade resources/deb-after-upgrade.sh \
		--before-remove resources/deb-before-remove.sh \
		build/linux-${GOARCH}/${BINARY}=/opt/policy-server-go/${BINARY} \
		resources/crossdomain.xml=/opt/policy-server-go/crossdomain.xml \
		resources/policy-server-monitor.sh=/opt/policy-server-go/policy-server-monitor.sh

fmt:
	go fmt

prep:
	mkdir -p build/darwin-${GOARCH}
	mkdir -p build/linux-${GOARCH}

clean:
	- rm -r build

.PHONY: linux darwin fmt clean prep package
