# app

![view](node-app.svg)

Application


<p align="right">
  [ <a href="../input/ndiag.descriptions/_node-app.md">:pencil2: Edit description</a> ]
</p>

## Components

| Name | Description | From (Relation) | To (Relation) |
| --- | --- | --- | --- |
| app:nginx |  <a href="../input/ndiag.descriptions/_component-app_nginx.md">:pencil2:</a> | [lb:nginx](node-lb.md) | [app:rails app](node-app.md) |
| app:rails app |  <a href="../input/ndiag.descriptions/_component-app_rails_app.md">:pencil2:</a> | [app:nginx](node-app.md) | [db:postgresql](node-db.md) |

## Labels

| Name | Description |
| --- | --- |
| [app](label-app.md) | App-DB |
| [http](label-http.md) | HTTP request flow |
| [lang:ruby](label-lang_ruby.md) | Ruby |
## Real nodes

- app-1
- app-2
- app-3

---

> Generated by [ndiag](https://github.com/k1LoW/ndiag)
