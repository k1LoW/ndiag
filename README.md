# ndiag

`ndiag` is a high-level architecture diagramming/documentation tool.

Key features of `ndiag` are:

- **N**ode-based: draw diagrams of nodes and their components. Nodes usually represent VMs, servers, or container hosts, but you can also redefine nodes.
- **N**ested-clusters: nodes can be clustered in layers.
- **N**-diagrams: generate multiple diagrams from a single configuration pair.

### Node

id = `[node name]`

### Layer

id = `[layer name]`

### Cluster

id = `[layer name]:[cluster name]`

### Component

**global component:**

id = `[component name]`

**cluster component:**

id = `[layer name]:[cluster name]:[component name]`

**node component:**

id = `[node name]:[component name]`
