# overview

![diagram](diagram-overview.svg)

全体イメージ


<p align="right">
  [ <a href="../ndiag.descriptions.ja/_diagram-overview.md">:pencil2: Edit description</a> ]
<p>



## 構成要素

| Name | Description |
| --- | --- |
| [real nodes](node-real_nodes.md) | 信頼できる情報源から得たNode ( [sample](/sample/input/nodes.yml#L1-L7) ) ... |
| [nodes](node-nodes.md) | Nodeをまとめるための設定 ( [sample](/sample/input/ndiag.yml#L13-L32) ) ... |
| [descriptions](node-descriptions.md) | `ndiag doc` によって生成したアーキテクチャドキュメントの各要素の説明文章 ( [sample](/sample/input/ndiag.descriptions) ) |
| [diagrams](node-diagrams.md) | [custom document](node-documents.md#components) を生成するための設定 ( [sample](/sample/input/ndiag.yml#L5-L12) ) ... |
| [networks](node-networks.md) | Component間のネットワーク ( [sample](/sample/input/ndiag.yml#L34-L59) ) ... |
| [relations](node-relations.md) | Component間の関係情報 ( [sample](/sample/input/ndiag.yml#L61-L67) ) |
| [ndiag](node-ndiag.md) | `ndiag doc` コマンド |
| [documents](node-documents.md) | 出力されるアーキテクチャドキュメント ( [sample](/sample/output/README.md) ) |


---

> Generated by [ndiag](https://github.com/k1LoW/ndiag)