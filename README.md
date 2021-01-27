# ndiag

`ndiag` is a "high-level architecture" diagramming/documentation tool.

Key features of `ndiag` are:

- **N**ode-based: draw diagrams of nodes and their components. Nodes usually represent VMs, servers, or container hosts, but you can also redefine nodes.
- **N**ested-clusters: nodes can be clustered in layers.
- **N**-diagrams: generate multiple diagrams from a single configuration.

## Usage

``` console
$ ndiag doc -c ndiag.yml --rm-dist
```

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
