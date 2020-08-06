# 記法変換サービス (TypeScript)
このディレクトリには記法変換 (Renderer) サービスの TypeScript による参考実装が含まれています.

主要なモジュール:

- `src/renderer.ts`: 記法変換の実装が含まれる
- `src/server.ts`: `src/renderer.ts` を gRPC で操作するためのインターフェース (gRPC サーバー) を実装する

その他のディレクトリ / モジュール:

- `pb/`: gRPC サービス定義から自動生成されたコード (リポジトリルートの `/pb/ts` からコピーされたもの)
- `src/config.ts`: サーバーの設定を読み込む

## テスト
以下のコマンドを実行します.

``` shell
yarn test
```
