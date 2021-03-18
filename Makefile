VERSION ?= "dev"
LDFLAGS=-ldflags "-s -w -X github.com/axllent/utproxy/cmd.Version=${VERSION}"
BINARY=utproxy

build = echo "\n\nBuilding $(1)-$(2)" && GO386=softfloat CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build ${LDFLAGS} -o dist/${BINARY}$(3) \
	&& if [ $(1) = "windows" ]; then \
		zip -9jq dist/${BINARY}_$(1)_$(2).zip dist/${BINARY}$(3) && rm -f dist/${BINARY}$(3); \
	else \
		tar --owner=root --group=root -C dist/ -zcf dist/${BINARY}_$(1)_$(2).tar.gz ${BINARY} && rm -f dist/${BINARY}; \
	fi

build: *.go go.*
	CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY}
	rm -rf /tmp/go-*

clean:
	rm -f ${BINARY}

release:
	mkdir -p dist
	rm -f dist/${BINARY}_*
	$(call build,linux,amd64,)
	$(call build,linux,386,)
	$(call build,darwin,amd64,)
	$(call build,darwin,arm64,)
	$(call build,windows,386,.exe)
	$(call build,windows,amd64,.exe)