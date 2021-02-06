# ndiag

`ndiag` is a "high-level architecture" diagramming/documentation tool.

Key features of `ndiag` are:

- **N**-diagrams: Generate multiple "diagrams and documents" (views) from a single configuration.
- **N**ested-clusters: Nodes and components can be clustered in nested layers
- Provides the ability to check the difference between the architecture and the actual system.

## Usage

``` console
$ ndiag doc -c ndiag.yml --rm-dist
```

## Getting Started

### Generate architecture documents

Add `ndiag.yml` (Full version is [here](example/tutorial/final/ndiag.yml)).

```yaml
---
name: 3-Tier Architecture
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - nginx?icon=lb-l7
    clusters:
      - 'vip_group:lb'
  -
    name: app
    components:
      - NGINX?icon=proxy
      - Rails App?icon=cube4&label=lang:ruby
  -
    name: db
    components:
      - PostgreSQL?icon=db

networks:
  -
    labels:
[...]
```

Run `ndiag doc` to generate documents in GitHub Friendly Markdown format.

``` console
$ ndiag doc -c ndiag.yml
```

Commit `ndiag.yml` and documents.

``` console
$ git add ndiag.yml doc/arch ndiag.description
$ git commit -m 'Add architecture document'
$ git push origin main
```

View the document on GitHub.

[Example document](example/tutorial/final/docs/arch/README.md)

![example](img/doc.png)

## Architecture

- [English](/docs/arch/README.md) :construction:
- [日本語](/docs/arch.ja/README.md)

## Install

**deb:**

Use [dpkg-i-from-url](https://github.com/k1LoW/dpkg-i-from-url)

``` console
$ export NDIAG_VERSION=X.X.X
$ curl -L https://git.io/dpkg-i-from-url | bash -s -- https://github.com/k1LoW/ndiag/releases/download/v$NDIAG_VERSION/ndiag_$NDIAG_VERSION-1_amd64.deb
```

**RPM:**

``` console
$ export NDIAG_VERSION=X.X.X
$ yum install https://github.com/k1LoW/ndiag/releases/download/v$NDIAG_VERSION/ndiag_$NDIAG_VERSION-1_amd64.rpm
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/ndiag
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/ndiag/releases)

**go get:**

```console
$ go get github.com/k1LoW/ndiag
```
