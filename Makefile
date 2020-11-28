PKG = github.com/k1LoW/ndiag
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test sec

test:
	go test ./... -coverprofile=coverage.txt -covermode=count

sec:
	gosec ./...

lint:
	golangci-lint run ./...

build:
	packr2
	go build -ldflags="$(BUILD_LDFLAGS)"
	packr2 clean

ndiag_draw: build
	./ndiag draw -c sample/input/ndiag.yml -n sample/input/nodes.yml -t png -l consul -l vip_group > sample.png
	./ndiag draw -c sample/input/ndiag.yml -n sample/input/nodes.yml -t svg -l consul -l vip_group > sample.svg
	./ndiag draw -c sample/input/ndiag.yml -n sample/input/nodes.yml -t dot -l consul -l vip_group > sample.dot
	./ndiag draw -c ndiag_ndiag.yml -n ndiag_ndiag.yml -t dot -l type -l file > ndiag.dot

ndiag_doc: build
	./ndiag doc -c ndiag_ndiag.yml -n ndiag_ndiag.yml --rm-dist
	./ndiag doc -c ndiag_ndiag.ja.yml -n ndiag_ndiag.ja.yml --rm-dist
	./ndiag doc -c sample/input/ndiag.yml -n sample/input/nodes.yml --rm-dist

ci_doc: depsdev ndiag_doc
	$(eval DIFF_EXIST := $(shell git checkout go.* && git diff --exit-code --quiet || echo "exist"))
	test -z "$(DIFF_EXIST)" || (git add -A ./docs && git add -A ./sample && git commit -m "Update ndiag archtecture document by GitHub Action (${GITHUB_SHA})" && git push -v origin ${GITHUB_BRANCH} && exit 1)

depsdev:
	go get github.com/Songmu/ghch/cmd/ghch
	go get github.com/gobuffalo/packr/v2/packr2
	go get github.com/Songmu/gocredits/cmd/gocredits
	go get github.com/securego/gosec/cmd/gosec

prerelease:
	git pull origin master --tag
	ghch -w -N ${VER}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS
	git commit -m'Bump up version number'
	git tag ${VER}

release:
	goreleaser --rm-dist

.PHONY: default test
