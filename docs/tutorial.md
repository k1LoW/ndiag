# Tutorial

In this tutorial, we will create a simple web service architecture documents.

To install `ndiag` command, please check the ["Install" section](../README.md#install).

## STEP1: Represent the roles of the instance and the middlewares/apps on the instance using "Node" and "Component"

**:pushpin: Keyword:** `Node`, `Component`, `Node component`

First, Create a YAML document as `ndiag.yml` like the following

```yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX
  -
    name: app
    components:
      - NGINX
      - App
  -
    name: db
    components:
      - PostgreSQL
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX
  -
    name: app
    components:
      - NGINX
      - App
  -
    name: db
    components:
      - PostgreSQL
```

[ndiag.yml in repo](../example/tutorial/step1/ndiag.yml)

</details>

In this `ndiag.yml`, the roles of the instances (`lb`, `app`, `db`) is represented by **Node** and the middlewares and applications on the instances are represented by **Component**.

Both Node and Component are elements that make up a system (architectural elements).

Then, run `ndiag doc` command.

``` console
$ ndiag doc -c ndiag.yml --rm-dist
```

If the command is successful, two directories should be created as follows.

``` console
$ ls
docs ndiag.descriptions ndiag.yml
$ tree docs/
docs/
└── arch
    ├── README.md
    ├── node-app.md
    ├── node-app.svg
    ├── node-db.md
    ├── node-db.svg
    ├── node-lb.md
    ├── node-lb.svg
    ├── view-nodes.md
    └── view-nodes.svg

1 directory, 9 files
$ tree ndiag.descriptions
ndiag.descriptions
├── _component-app_app.md
├── _component-app_nginx.md
├── _component-db_postgresql.md
├── _component-lb_nginx.md
├── _index.md
├── _node-app.md
├── _node-db.md
├── _node-lb.md
└── _view-nodes.md

0 directories, 9 files
```

| Directory | |
| --- | --- |
| `docs/` | Generated documents |
| `ndiad.descriptions` | Sub documents to set description of architecture elements ( It will be explained in STEP7 ) |

Open the file `docs/arch/README.md`. The documentation template is now complete.

### :book: Generated documents of this STEP

<img src="../example/tutorial/step1/docs/arch/view-nodes.svg" />

[Generated documents](../example/tutorial/step1/docs/arch/README.md)

### :memo: Point of this STEP

A Component that belongs to Node is called Node component.

## STEP2: Represent data flow (HTTP request/Database access etc) using "networks:"

**:pushpin: Keyword:** `networks:`, `Global component`

Data flow between components (HTTP requests/database access, etc.) is represented by adding `networks:` to `ndiag.yml`.

``` yaml
[...]

networks:
  -
    route:
      - "internet"
      - "vip"
  -
    route:
      - "vip"
      - "lb:nginx"
  -
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
[...]
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX
  -
    name: app
    components:
      - NGINX
      - App
  -
    name: db
    components:
      - PostgreSQL

networks:
  -
    route:
      - "internet"
      - "vip"
  -
    route:
      - "vip"
      - "lb:nginx"
  -
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
  -
    route:
      - "app:app"
      - "Payment API"
  -
    route:
      - "app:app"
      - "db:postgresql"
```

[ndiag.yml in repo](../example/tutorial/step2/ndiag.yml)

</details>

Then, run `ndiag doc` command as in STEP1.

``` console
$ ndiag doc -c ndiag.yml --rm-dist
```

( After STEP2, execute `ndiag doc` command to generate the document. )

### :book: Generated documents of this STEP

<img src="../example/tutorial/step2/docs/arch/view-nodes.svg" />

[Generated documents](../example/tutorial/step2/docs/arch/README.md)

### :memo: Point of this STEP

Node component is specified by joining "Node id (= Node name)" and "Component name" with `:`.

**:bulb: Example:** `lb:nginx` means "Component `NGINX`" that belongs to "Node `lb`"

A Component that does not belong to Node (or Cluster) is called "Global component" ( `internet`, `vip`, `Payment API` ). It is specified by only the Component name.

## STEP3: Represent relationships between components other than the data flow using "relations:"

**:pushpin: Keyword:** `relations:`

The relations between Components, other than the data flow, are expressed using `relations:` as shown below.

``` yaml
[...]

relations:
  -
    components:
      - 'lb:Keepalived'
      - "vip"
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX
  -
    name: app
    components:
      - NGINX
      - App
  -
    name: db
    components:
      - PostgreSQL

networks:
  -
    route:
      - "internet"
      - "vip"
  -
    route:
      - "vip"
      - "lb:nginx"
  -
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
  -
    route:
      - "app:app"
      - "db:postgresql"
  -
    route:
      - "app:app"
      - "Payment API"

relations:
  -
    components:
      - 'lb:Keepalived'
      - "vip"
```

[ndiag.yml in repo](../example/tutorial/step3/ndiag.yml)

</details>

Then, run `ndiag doc` command.

``` console
$ ndiag doc -c ndiag.yml --rm-dist
```

### :book: Generated documents of this STEP

<img src="../example/tutorial/step3/docs/arch/view-nodes.svg" />

[Generated documents](../example/tutorial/step3/docs/arch/README.md)

### :memo: Point of this STEP

`networks:` is another expression for `type: network` in `relations:`. You can use either one.

**:bulb: Example:**

<table>
  <tr><th> networks: </th><th> relations: </th></tr>
  <tr>
    <td>
<pre>
networks:
  -
    route:
      - "internet"
      - "vip"
</pre>
    </td>
    <td>
<pre>
relations:
  -
    type: network
    components:
      - "internet"
      - "vip"
</pre>
    </td>
  </tr>
</table>

## STEP4: Grouping nodes and components using "Cluster" and "Layer"

**:pushpin: Keyword:** `Cluster`, `Layer`, `Cluster component`

Represent groups of Nodes and Components using `clusters:`.

``` yaml
[...]
  -
    name: db
    components:
      - PostgreSQL
    clusters:
      - 'consul:dc1'

networks:
  -
    route:
      - "internet"
      - "vip_group:lb:vip"
[...]
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX
    clusters:
      - 'Consul:dc1'
      - 'vip_group:lb'
  -
    name: app
    components:
      - NGINX
      - App
    clusters:
      - 'consul:dc1'
  -
    name: db
    components:
      - PostgreSQL
    clusters:
      - 'consul:dc1'

networks:
  -
    route:
      - "internet"
      - "vip_group:lb:vip"
  -
    route:
      - "vip_group:lb:vip"
      - "lb:nginx"
  -
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
  -
    route:
      - "app:app"
      - "db:postgresql"
  -
    route:
      - "app:app"
      - "Service:Payment:Payment API"

relations:
  -
    components:
      - 'lb:Keepalived'
      - "vip_group:lb:vip"
```

[ndiag.yml in repo](../example/tutorial/step4/ndiag.yml)

</details>

"Node" can belong to multiple Clusters.

**:bulb: Example:**

``` yaml
[...]
nodes:
  -
    name: instance
    components:
      - http-server
    clusters:
      - 'role:web'
      - 'location:dc'
      - 'os:ubuntu-focal'
[...]
```

### :book: Generated documents of this STEP

<img src="../example/tutorial/step4/docs/arch/view-nodes.svg" />

[Generated documents](../example/tutorial/step4/docs/arch/README.md)

### :memo: Point of this STEP

In ndiag, Nodes and Components can be grouped by an element called **Cluster**.

A Node can belong to multiple Clusters.

**:bulb: Example:**

``` yaml
[...]
nodes:
  -
    name: instance
    components:
      - http-server
    clusters:
      - 'role:web'
      - 'location:dc'
      - 'os:ubuntu-focal'
[...]
```

"Cluster" always belongs to a "Layer". "Layer" can have multiple Clusters.

In the figure, Clusters that belong to the same Layer are represented by lines of the same color.

"Cluster" is specified by joining Layer id (= Layer name) and Cluster name with `:`.

**:bulb: Example:** Layer `role` has a Cluster `role:web` and a Cluster `role:db` with the same Layer id `role`.

Also, a Component that belongs to Cluster instead of Node is called **Cluster component**.

Cluster component is specified by joining Cluster id and Component name with `:`.

**:bulb: Example:** `vip_group:lb:vip` means "Component `vip`" that belongs to "Cluster `vip_group:lb`"

## STEP5: Add icons

**:pushpin: Keyword:** `icon`

:construction:

``` yaml
[...]
  -
    name: db
    components:
      - PostgreSQL?icon=db
    clusters:
      - 'consul:dc1'

networks:
  -
    route:
      - "internet?icon=cloud"
      - "vip_group:lb:vip"
[...]
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
nodes:
  -
    name: lb
    components:
      - NGINX?icon=lb-l7
    clusters:
      - 'Consul:dc1'
      - 'vip_group:lb'
  -
    name: app
    components:
      - NGINX?icon=proxy
      - App?icon=cube4
    clusters:
      - 'consul:dc1'
  -
    name: db
    components:
      - PostgreSQL?icon=db
    clusters:
      - 'consul:dc1'

networks:
  -
    route:
      - "internet?icon=cloud"
      - "vip_group:lb:vip"
  -
    route:
      - "vip_group:lb:vip"
      - "lb:nginx"
  -
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
  -
    route:
      - "app:app"
      - "db:postgresql"
  -
    route:
      - "app:app"
      - "Service:Payment:Payment API"

relations:
  -
    components:
      - 'lb:Keepalived?icon=keepalived'
      - "vip_group:lb:vip"

customIcons:
  -
    key: keepalived
    lines:
      - b1 b5 f9 j5 j1 f1 b1
      - d2 d6
      - h2 d4
      - e4 h6
```

[ndiag.yml in repo](../example/tutorial/step5/ndiag.yml)

</details>

### :book: Generated documents of this STEP

<img src="../example/tutorial/step5/docs/arch/view-nodes.svg" />

[Generated documents](../example/tutorial/step5/docs/arch/README.md)

## STEP6: Create architecture views using "Label" and "views:"

**:pushpin: Keyword:** `views:`, `Label`

:construction:

``` yaml
[...]
views:
  -
    name: overview
    layers: ["consul", "vip_group", "service"]
[...]
```
<details>

<summary> Full version of <code>ndiag.yml</code> is here (click).</summary>

``` yaml
---
name: Simple web service
docPath: docs/arch
views:
  -
    name: overview
    layers: ["consul", "vip_group", "service"]
  -
    name: http access
    layers: ["vip_group"]
    labels: ["http"]
  -
    name: app
    layers: ["vip_group", "service"]
    labels: ["app"]
nodes:
  -
    name: lb
    components:
      - NGINX?icon=lb-l7
    clusters:
      - 'Consul:dc1'
      - 'vip_group:lb'
  -
    name: app
    components:
      - NGINX?icon=proxy
      - App?icon=cube4&label=lang:ruby
    clusters:
      - 'consul:dc1'
  -
    name: db
    components:
      - PostgreSQL?icon=db
    clusters:
      - 'consul:dc1'

networks:
  -
    labels:
      - http
    route:
      - "internet?icon=cloud"
      - "vip_group:lb:vip"
  -
    labels:
      - http
    route:
      - "vip_group:lb:vip"
      - "lb:nginx"
  -
    labels:
      - http
      - app
    route:
      - "lb:nginx"
      - "app:nginx"
      - "app:app"
  -
    labels:
      - app
    route:
      - "app:app"
      - "db:postgresql"
  -
    labels:
      - app
    route:
      - "app:app"
      - "Service:Payment:Payment API"

relations:
  -
    labels:
      - http
    components:
      - 'lb:Keepalived?icon=keepalived'
      - "vip_group:lb:vip"

customIcons:
  -
    key: keepalived
    lines:
      - b1 b5 f9 j5 j1 f1 b1
      - d2 d6
      - h2 d4
      - e4 h6
```

[ndiag.yml in repo](../example/tutorial/step6/ndiag.yml)

</details>

### :book: Generated documents of this STEP

<img src="../example/tutorial/step6/docs/arch/view-overview.svg" />

[Generated documents](../example/tutorial/step6/docs/arch/README.md)

## STEP7: Add descriptions using GitHub Web UI and GitHub Actions

:construction:

### :book: Generated documents of this STEP

[Generated documents](../example/tutorial/step7/docs/arch/README.md)

## STEP8: Strict mode

:construction:
