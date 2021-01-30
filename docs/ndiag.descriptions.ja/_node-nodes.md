Nodeをまとめるための設定 ( [example](/example/input/ndiag.yml#L13-L32) )

nodes.ymlに以下のように3つのNodeがあり、

``` yml
- app-1
- app-2
- app-3
```

ndiag.ymlにワイルドカードを含んだ `app-*` がある場合、アーキテクチャドキュメントでは `app-1` `app-2` `app-3` は `app-*` にまとめられて表現されます。

これは、スケールアウトさせている同じ構成のインスタンスをドキュメントでわかりやすく省略表現する場合などに利用します。

``` yml
nodes:
  -
    name: app-*
    components:
      - rails
    clusters:
      - 'consul:dc1'
```
