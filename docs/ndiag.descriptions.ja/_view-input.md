## 入力

`ndiag doc` の入力は以下の3種類です。

1. [nodes.yml](/example/3-tier/input/nodes.yml)
    - 信頼できる情報源から得たNodeやComponent、Network、Relationのデータです。
    - 信頼できる情報源とするために、例えばあるインフラのアーキテクチャドキュメントを生成する場合、インフラから収集したサーバやVMインスタンスの情報やミドルウェアの情報からnodes.ymlを機械的に生成することが望ましいです。
        - Terraformの `*.tf` や `*.tfstate`
        - Chefの `nodes/`
        - `kubectl get nodes` `kubectl get pods`
        - etc.
    - できる限り最新の情報を継続的に取得できるような仕組みを作ると良いでしょう。
2. [ndiag.yml](/example/3-tier/input/ndiag.yml)
    - 生成するドキュメントの設定と、nodes.ymlとして生成できなかった情報を補完するための設定ファイルです。
        - ドキュメントの出力ディレクトリ `docPath:`
        - 説明文章サブセットmdファイルのディレクトリ `descPath:`
        - ドキュメントの体裁の設定
        - Nodeをまとめるための設定 ( `nodes:` )
        - etc.
3. [ndiag.descriptions/*.md (descPath)](/example/3-tier/input/ndiag.descriptions)
    - NodeやComponent、生成したアーキテクチャ図の説明文章を記載するためのサブセットmdファイルです。
    - `ndiag doc` を実行すると、ndiag.ymlの `descPath` に設定されたディレクトリにサブセットmdファイルが（存在しない場合）空ファイルとして生成されます。

入力処理の流れとしては以下のようになります。

1. nodes.ymlのNodeの情報をndiag.ymlの `nodes:` の情報でまとめる
2. nodes.ymlのNode以外の情報をndiag.ymlの `components:` `relations:` `networks:` などにマージしていく
3. 2まででまとめた情報から必要な説明文章サブセットmdファイルを `descPath:` ( default `ndiag.descriptions/` ) に指定されたディレクトリからファイル名を指定して読み込む。
    - もしファイルがなければ空ファイルを生成する

３の時点で `ndiag:intenal` には必要な情報が構造化されて保持されています。
