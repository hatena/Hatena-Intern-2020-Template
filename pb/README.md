# gRPC サービスの定義
このディレクトリでは認証基盤と記法変換の 2 つの gRPC サービスの定義を行います.

- `account.proto`, `renderer.proto` には, それぞれ gRPC サービスの定義が [protocol buffers](https://developers.google.com/protocol-buffers/docs/proto3) で記述されています
- これらの定義から自動生成されたコードは `go`, `ts` の各ディレクトリに配置され, また各サービスの実装のディレクトリにもコピーされています

## コード生成
`scripts/compile` を実行すると

1. コード生成を行うための Docker イメージの作成
2. コード生成
3. 生成したコードの各サービスの実装ディレクトリへのコピー

を行います.

