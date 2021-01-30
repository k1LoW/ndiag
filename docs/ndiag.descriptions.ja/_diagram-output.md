## 出力 (アーキテクチャドキュメント)

`ndiag doc` が出力するドキュメントは ndiag.yml の `docPath` に設定されたディレクトリ（デフォルトは `archdoc` ）に生成されます。

- [output/README.md (docPath)](/example/3-tier/output/README.md)

ドキュメントは1つではなく、複数生成します。

### index document

[output/README.md (docPath)](/example/3-tier/output/README.md)

### layer based document

ndiag.ymlで設定したLayerごとにドキュメントを生成します。

それぞれのLayerを中心とした説明をすることに使用します。

- [output/layer-consul.md](/example/3-tier/output/layer-consul.md)

### label based document

ndiag.ymlで設定したrelationsやnetworksに付与したLabelごとにドキュメントを生成します。

Labelで表したComponentの関係を中心とした説明をすることに使用します。

- [output/label-http.md](/example/3-tier/output/label-http.md)

### custom document

ndiag.ymlのdiagramsで設定したlayers、labelsを元にドキュメントを生成します。

- [output/diagram-http-lb.md](/example/3-tier/output/diagram-http-lb.md)
