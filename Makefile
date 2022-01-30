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
	./ndiag draw -c example/3-tier/input/ndiag.yml -n example/3-tier/input/nodes.yml -t png -l consul -l vip_group > example.png
	./ndiag draw -c example/3-tier/input/ndiag.yml -n example/3-tier/input/nodes.yml -t svg -l consul -l vip_group > example.svg
	./ndiag draw -c example/3-tier/input/ndiag.yml -n example/3-tier/input/nodes.yml -t dot -l consul -l vip_group > example.dot
	./ndiag draw -c ndiag_ndiag.yml -t dot -l type -l file > ndiag.dot

ndiag_doc: ndiag_doc_local ndiag_doc_tutorial ndiag_doc_fetch_icons

ndiag_doc_local: build
	./ndiag doc -c ndiag_ndiag.yml --rm-dist
	./ndiag doc -c ndiag_ndiag.ja.yml --rm-dist
	./ndiag doc -c example/3-tier/input/ndiag.yml -n example/3-tier/input/nodes.yml --rm-dist

ndiag_doc_tutorial: build
	./ndiag doc -c example/tutorial/step1/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step2/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step3/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step4/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step5/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step6/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/step7/ndiag.yml --rm-dist
	./ndiag doc -c example/tutorial/final/ndiag.yml -n example/tutorial/final/nodes.yml --rm-dist

ndiag_doc_fetch_icons: build
	./ndiag fetch-icons k8s -c example/k8s/input/ndiag.yml && touch example/k8s/input/ndiag.icons/.gitkeep && echo "*.*" > example/k8s/input/ndiag.icons/.gitignore
	./ndiag doc -c example/k8s/input/ndiag.yml --rm-dist
	./ndiag fetch-icons gcp -c example/gcp/input/ndiag.yml && touch example/gcp/input/ndiag.icons/.gitkeep && echo "*.*" > example/gcp/input/ndiag.icons/.gitignore
	./ndiag doc -c example/gcp/input/ndiag.yml --rm-dist
	./ndiag fetch-icons aws -c example/aws/input/ndiag.yml && touch example/aws/input/ndiag.icons/.gitkeep && echo "*.*" > example/aws/input/ndiag.icons/.gitignore
	./ndiag doc -c example/aws/input/ndiag.yml --rm-dist

ci_doc: depsdev ndiag_doc

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@v0.10.2
	go install github.com/Songmu/gocredits/cmd/gocredits@v0.2.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.8.1
	go install github.com/gobuffalo/packr/v2/packr2@v2.8.3

prerelease:
	git pull origin main --tag
	go mod tidy
	ghch -w -N ${VER}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

release:
	goreleaser --rm-dist

.PHONY: default test
