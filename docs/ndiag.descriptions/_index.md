This document explains how ndiag generates the architecture documentation when you actually run the following commands.

#### Sample Architecture

WIP

[3-Tier Architecture](/example/3-tier/output/README.md)

![3tier](/example/3-tier/output/diagram-overview.svg)

``` console
$ ndiag doc -c example/3-tier/input/ndiag.yml -n example/3-tier/input/nodes.yml --rm-dist
```

### Documentation cycle

WIP

#### Input

WIP

#### Output (architecture document)

WIP

### Node

node id = `[node name]`

### Layer

layer id = `[layer name]`

### Cluster

cluster id = `[layer name]:[cluster name]`

### Component

**global component:**

component id = `[component name]`

**cluster component:**

component id = `[layer name]:[cluster name]:[component name]`

**node component:**

component id = `[node name]:[component name]`
