# ndiag

`ndiag` is a high-level architecture diagramming/documentation tool.

Key features of `ndiag` are:

- **N**ode-based: draw diagrams of nodes and their components. Nodes usually represent VMs, servers, or container hosts, but you can also redefine nodes.
- **N**ested-clusters: nodes can be clustered in layers.
- **N**-diagrams: generate multiple diagrams from a single configuration pair.

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
