[custom document](node-documents.md#components) を生成するための設定 ( [sample](/sample/input/ndiag.yml#L5-L12) )

サンプルアーキテクチャのndiag.ymlの `diagrams:` の設定は以下のようになっています。

``` yaml
diagrams:
  -
    name: overview
    layers: ["consul", "vip_group"]
  -
    name: http-lb
    layers: ["vip_group"]
    labels: ["http"]
```

`layers:` にClusterがうまく入子構造になるようにLayerを設定することで意図したグルーピングでNodeやComponentを配置した図とドキュメントの雛形を作成できます。

また、`labels:` でLabelを（複数）指定することで、Labelで関係性を持ったコンポーネントだけに限定した図とドキュメントの雛形を作成できます。
