# 課題

Hatena SUMMER INTERNSHIP 2020では、用意されたブログシステムに追加の機能を実装してもらいます。ブログシステムは複数のサービスから成りますが、特に「記法変換サービス」を実装するのが、皆さんの課題です。

## 準備

### リポジトリの用意

[hatena/Hatena-Intern-2020-Template](https://github.com/hatena/Hatena-Intern-2020-Template)に、課題に取り組むために必要なコードがあります。テンプレートリポジトリになっているので、「Use Template」ボタンを押して、自分のリポジトリを作ってください。

### セットアップ

リポジトリのREADMEに従って、手持ちのコンピュータに環境を構築してください。コンピュータは、以下を満たすものを用意してください。

|  | 推奨環境 |
|:--|:--|
| OS | Windows/macOS/Linux |
| CPU | 64-bitのマルチコアCPU |
| メモリ | 16 GB以上 |
| ストレージ | SSD |

## 課題

ブログシステムのコードはGo言語で書かれています。課題はGo言語もしくはTypeScriptで取り組むことを想定しています。しかしgRPCでやりとりできるなら、他の任意の言語で実装しても構いません。

### 課題1: 記法変換サービス

ブログの記事をMarkdownで書けたら嬉しいと思います。そこで記事本文をなんらかの「記法」で装飾できるようにしてください。このとき、少なくとも以下の3つの記法を実装してください。

- 見出し記法
- リスト記法
- リンク記法

記法は[Markdown](https://commonmark.org/help/)でも、独自に作成しても構いません。Markdownを採用する場合、ライブラリを利用してもよいです。

参考までに、Go言語でライブラリを利用する場合は[goldmark](https://github.com/yuin/goldmark)を、TypeScriptでは[unified](https://github.com/unifiedjs/unified)（[remark](https://github.com/remarkjs/remark)、[rehype](https://github.com/rehypejs/rehype)）を推奨します。

テンプレートの中の[renderer-go](https://github.com/hatena/Hatena-Intern-2020-Template/tree/master/services/renderer-go)サービス、もしくは[renderer-ts](https://github.com/hatena/Hatena-Intern-2020-Template/tree/master/services/renderer-ts)サービスが出発点になるでしょう。

#### 発展

独自の記法を考えて追加してみましょう。自分がほしいと思えるような記法であると、なおよいです。

Markdownのライブラリを利用している場合も、上で推奨したライブラリであれば、拡張可能になっています。ドキュメントをよく読んで、うまく拡張してください。

### 課題2: タイトルの自動取得サービス

ブログの記事にリンクを載せたいときに、リンク先のタイトルを自動的に埋めることができたら、便利ではないでしょうか。リンク記法を拡張して、リンクテキストを省略した場合に、自動的にページタイトルが使われるようにしてください。Markdownであれば以下のような入出力になります。

入力: `[](https://example.com)`

出力: `<a href="https://example.com">Example Domain</a>`

ページタイトルを取得するのは、記法を変換することとは異なったビジネスドメインなので、独立したサービスとして実装してください。

記法変換サービスを`renderer`、タイトル取得サービスを`fetcher`とすると、実装は大まかに、以下のような手順で行うとよいでしょう。

- `services/`以下に新しく`fetcher`サービスを作る
    - `skaffold.yml`や`k8s/`以下のマニフェストを編集してサービスを起動させる
- `pb/fetcher.proto`を作り、`scripts/compile`を変更して、サーバー・クライアントのコードが生成されるようにする
- `fetcher`サービスに、ページタイトルを取得する実装を行う
- `renderer`サービスから適宜`fetcher`サービスを呼び出し、リンクテキストを適切に設定する

#### 発展

- URLが20個並んだブログ記事を更新できるようにするために、どうすればいいいか考えて実装してください。
    - URLが100個、あるいは500個あるときはどうすればいいか、考えてみてください。
- タイトル取得サービスでは、外部のウェブサイトへリクエストを行います。このようなサービスを設計する際に、注意すべきことを考えましょう。
    - ウェブサイト運営者の立場に立って、困ることを考えましょう。
    - クローラーを実装するときに、どのような対応が必要でしょうか。
    - 例：リクエスト頻度、robots.txt、クロール禁止サイト（[Robots Beware](https://arxiv.org/help/robots)）

## 補足

### テストについて

なるべくテストコードを書いてみましょう。

テストの目的の一つは、ソフトウェアの品質を保証することです。テストコードは、実装が正しいことを保証します。また、バグを修正した後に、修正されていることを確認するテストを書いておくことで、同じバグが再び起きていないことを将来にわたって保証することも可能です。

単に品質を保証するだけでなく、テストコードによって、テストの対象となるソフトウェアがどのように使われることを想定しているのか、示すこともできます。このようなテストコードはドキュメンテーションテストとして、ドキュメントの一部として扱われることもあります。

テストコードが果たすもう一つの役割は、生産性の向上です。ソフトウェアの動作を確認するのを、毎回人力でやっていると、手間がかかります。テストコードで自動化できると、とても便利です。ライブラリのアップデートのように、影響が広範にわたることが想定される場合、特に顕著です。

課題では、システムが複数のサービスに分けられています。全てのサービスを組み合わせないと動作が確認できない、となると大変です。個々のサービスを個別にテストできるように工夫してみましょう。

### コミットの大きさについて

コミットの粒度を意識して取り組んでみましょう。

全て実装してからまとめてコミットするようなやり方だと、gitのうまみが失われてしまいます。ひとつのコミットにどのような意味を持たせるのか考えて、適切なまとまりでコミットしてください。

ただしコミットログを気にして何度もrebaseするようなことは求めません。

## 参考

goldmarkを利用して、Markdownのリンク記法でリンクテキストが指定されていない際に、タイトルを差し込むコードの例

```golang
package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var markdown = goldmark.New(
	goldmark.WithParserOptions(
		parser.WithASTTransformers(
			util.Prioritized(&autoTitleLinker{}, 999),
		),
	),
)

func main() {
	src := []byte("# link samples\n" +
	"[normal link](https://example.com)\n" +
	"[](https://example.com)\n")
	var buf bytes.Buffer
	if err := markdown.Convert(src, &buf); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf.String())
}

type autoTitleLinker struct {
	// fetcherCli pb.FetcherClient (ヒント)
}

func (l *autoTitleLinker) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			node.AppendChild(node, ast.NewString([]byte(fetchTitle(string(node.Destination)))))
		}
		return ast.WalkContinue, nil
	})
}

func fetchTitle(url string) string {
	return "example title"
}
```
