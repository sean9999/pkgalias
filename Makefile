REPO=github.com/sean9999/pkgalias
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

binaries: bin/pkgalias
	mkdir -p bin

bin/pkgalias:
	go build -v -o bin/pkgalias -ldflags="-X 'main.Version=$(REF)'" cmd/pkgalias/main.go

tidy:
	go mod tidy

install:
	go install ./cmd/pkgalias

clean:
	go clean
	rm bin/*

pkgsite:
	if [ -z "$$(command -v pkgsite)" ]; then go install golang.org/x/pkgsite/cmd/pkgsite@latest; fi

docs: pkgsite
	pkgsite -open .

publish:
	GOPROXY=https://proxy.golang.org,direct go list -m ${REPO}@${SEMVER}

test:
	go test ./...

.PHONY: test
