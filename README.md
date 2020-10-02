# [Hatena SUMMER INTERNSHIP 2020](https://hatenacorp.jp/intern2020)

「Hatena SUMMER INTERNSHIP 2020」では、Kubernetes上に構築されたブログシステムを題材としました。ブログシステムはマイクロサービスを意識しており、メインであるブログサービスに加えて、アカウントサービスや、Markdownなどの記法を変換するサービスが用意されています。それぞれのサービス間はgRPCを使ってやりとりしています。

インターンシップのカリキュラムについては、[講義動画](https://hatenacorp.jp/intern2020/public_broadcast)や[課題](/docs/exercise.md)を公開しているので、参照してください。

## セットアップ
アプリケーションの起動には以下が必要です.

- [Docker](https://docs.docker.com/engine/install/)
  - Windows または macOS の場合は Docker Desktop
  - Linux の場合は各ディストリビューションごとのインストール方法に従ってください
- Kubernetes
  - 以下のいずれかを利用することを想定しています
    - Docker Desktop に付属のもの
    - [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kustomize](https://kubernetes-sigs.github.io/kustomize/installation/)
- [Skaffold](https://skaffold.dev/docs/install/)

個々のサービスの開発には, 以下がローカル環境にインストールされていることを想定しています.

- Make
- (Go を使う場合) [Go](https://golang.org/)
- (TypeScript を使う場合) [Node.js](https://nodejs.org/en/), [Yarn](https://classic.yarnpkg.com/lang/en/)

動作確認は以下の環境で行っています.

- macOS 10.15.6 (19G73)
- Docker Desktop 2.3.0.4 (46911) (Kubernetes v1.16.5)

``` console
$ docker version
Client: Docker Engine - Community
 Version:           19.03.12
 API version:       1.40
 Go version:        go1.13.10
 Git commit:        48a66213fe
 Built:             Mon Jun 22 15:41:33 2020
 OS/Arch:           darwin/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          19.03.12
  API version:      1.40 (minimum version 1.12)
  Go version:       go1.13.10
  Git commit:       48a66213fe
  Built:            Mon Jun 22 15:49:27 2020
  OS/Arch:          linux/amd64
  Experimental:     true
 containerd:
  Version:          v1.2.13
  GitCommit:        7ad184331fa3e55e52b890ea95e65ba581ae3429
 runc:
  Version:          1.0.0-rc10
  GitCommit:        dc9208a3303feef5b3839f4323d9beb36df0a9dd
 docker-init:
  Version:          0.18.0
  GitCommit:        fec3683

$ minikube version
minikube version: v1.12.2
commit: be7c19d391302656d27f1f213657d925c4e1cfc2

$ kubectl version --client
Client Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.6", GitCommit:"dff82dc0de47299ab66c83c626e08b245ab19037", GitTreeState:"clean", BuildDate:"2020-07-16T00:04:31Z", GoVersion:"go1.14.4", Compiler:"gc", Platform:"darwin/amd64"}

$ kustomize version
{Version:3.8.0 GitCommit:6a50372dd5686df22750b0c729adaf369fbf193c BuildDate:2020-07-05T17:55:53+01:00 GoOs:darwin GoArch:amd64}

$ skaffold version
v1.13.0

$ go version
go version go1.14.6 darwin/amd64

$ node -v
v14.6.0

$ yarn -v
1.22.4
```

## 起動
### Docker Desktop
以下の手順でアプリケーションを起動します.

``` shell
# context を設定
kubectl config set-context hatena-intern-2020 --cluster=docker-desktop --user=docker-desktop --namespace=hatena-intern-2020
kubectl config use-context hatena-intern-2020

# 起動
make up
```

ブラウザで http://localhost:8080/ を開くことでアプリケーションにアクセスできます.

### Minikube
以下の手順でアプリケーションを起動します.

``` shell
# Minikube を起動
minikube start --kubernetes-version v1.16.13
eval $(minikube docker-env)

# context を設定
kubectl config set-context hatena-intern-2020 --cluster=minikube --user=minikube --namespace=hatena-intern-2020
kubectl config use-context hatena-intern-2020

# 起動
make up
```

以下のコマンドを実行するとブラウザが自動的に開き, アプリケーションにアクセスします.

``` shell
minikube -n hatena-intern-2020 service blog
```

## サービス
アプリケーションには以下の 3 つのサービスが存在します.

- 認証基盤 (Account) サービス
  - ユーザーアカウントの登録や認証を管轄します
- ブログ (Blog) サービス
  - ユーザーに対して, ブログを作成したり記事を書いたりする機能を提供します
- 記法変換 (Renderer) サービス
  - ブログの記事を記述するための「記法」から HTML への変換を担います

このうちブログサービスが Web サーバーとして動作し, ユーザーに対してアプリケーションを操作するためのインターフェースを提供します.
認証基盤サービスと記法変換サービスは gRPC サービスとして動作し, ブログサービスから使用されます.

## ディレクトリ構成

- `pb/`: gRPC サービスの定義
- `services/`: 各サービスの実装
  - `account/`: 認証基盤サービス
  - `blog/`: ブログサービス
  - `renderer-go/`: 記法変換サービスの Go による実装
  - `renderer-ts/`: 記法変換サービスの TypeScript による実装
- `k8s/`: アプリケーションを Kubernetes 上で動作させるためのマニフェスト

## クレジット
- 株式会社はてな
  - [@akiym](https://github.com/akiym)
  - [@cockscomb](https://github.com/cockscomb)
  - [@itchyny](https://github.com/itchyny)
  - [@susisu](https://github.com/susisu)

(順不同)

このリポジトリの内容は MIT ライセンスで提供されます. 詳しくは `LICENSE` をご確認ください.
