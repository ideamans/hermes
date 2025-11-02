# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Hermes - Transactional Email Generator

Hermesは、トランザクショナルメール（ウェルカムメール、パスワードリセット、レシートなど）のクリーンでレスポンシブなHTML/プレーンテキストを生成するGoパッケージです。Node.jsの[mailgen](https://github.com/eladnava/mailgen)のGoポート版です。

### 開発コマンド

```bash
# パッケージのビルド
go build

# テストの実行
go test ./...

# 特定のテストのみ実行
go test -run TestThemeSimple

# カバレッジ付きテスト
go test -cover ./...

# サンプルの生成（examples/フォルダで実行）
cd examples
go run *.go

# サンプルを実際のメールアドレスに送信してテスト（環境変数設定が必要）
HERMES_SEND_EMAILS=true \
HERMES_SMTP_SERVER=smtp.gmail.com \
HERMES_SMTP_PORT=465 \
HERMES_SENDER_EMAIL=your@gmail.com \
HERMES_SENDER_IDENTITY="Your Name" \
HERMES_SMTP_USER=your@gmail.com \
HERMES_TO=recipient@example.com \
go run *.go
```

### アーキテクチャと技術スタック

**依存関係:**
- Go 1.24.2+
- `github.com/Masterminds/sprig` - テンプレート関数
- `github.com/russross/blackfriday/v2` - Markdown → HTML変換
- `github.com/vanng822/go-premailer` - CSS インライン化
- `github.com/jaytaylor/html2text` - HTML → プレーンテキスト変換
- `github.com/stretchr/testify` - テスト用アサーション

### コアアーキテクチャ

**主要ファイル:**
- `hermes.go` - メインジェネレーター。`Hermes`構造体とテンプレート生成ロジック
- `default.go` - Defaultテーマ実装（Postmark Transactional Email Templates）
- `flat.go` - Flatテーマ実装（Postmarkの改変版）
- `hermes_test.go` - 全テーマの包括的なテストスイート

**テーマシステム:**

新しいテーマを作成する場合は`Theme`インターフェースを実装:
```go
type Theme interface {
    Name() string              // テーマ名
    HTMLTemplate() string      // HTML用Goテンプレート
    PlainTextTemplate() string // プレーンテキスト用Goテンプレート
}
```

テーマは以下のような構成で実装:
1. テーマ構造体を定義（例: `type Default struct {}`）
2. `Name()`, `HTMLTemplate()`, `PlainTextTemplate()` メソッドを実装
3. `hermes_test.go`の`testedThemes`スライスに追加してテスト対象に含める

**メール生成フロー:**

1. `Hermes`インスタンスを作成（テーマ、商品情報を設定）
2. `Email`構造体でメールコンテンツを定義（本文、アクション、テーブルなど）
3. `GenerateHTML()`または`GeneratePlainText()`を呼び出し
4. 内部処理:
   - デフォルト値のマージ（`setDefaultHermesValues`, `setDefaultEmailValues`）
   - Goテンプレートの実行（sprig関数、カスタム関数を含む）
   - HTML生成の場合はCSS インライン化（Premailer使用、`DisableCSSInlining`で無効化可能）
   - プレーンテキスト生成の場合はHTML→テキスト変換（html2text使用）

### テスト設計

**テストパターン:**

全テーマが同じ情報を表示することを保証するため、`Example`インターフェースを使用:
```go
type Example interface {
    getExample() (h Hermes, email Email)           // テストデータの作成
    assertHTMLContent(t *testing.T, s string)      // HTML出力の検証
    assertPlainTextContent(t *testing.T, s string) // プレーンテキスト出力の検証
}
```

各機能ごとに`Example`実装を作成（例: `SimpleExample`, `WithInviteCode`, `WithFreeMarkdownContent`）し、`testedThemes`の全テーマで実行。

**新しいテーマを追加する際:**
1. テーマ構造体と`Theme`インターフェースメソッドを実装
2. `hermes_test.go`の`testedThemes`スライスに`new(YourTheme)`を追加
3. `go test ./...`を実行して全既存テストがパスすることを確認

### 重要な実装ポイント

**CSSインライン化:**
デフォルトでは`go-premailer`を使用してCSSをインライン化し、メールクライアントとの互換性を向上。`Hermes.DisableCSSInlining = true`で無効化可能。

**マージ戦略:**
`dario.cat/mergo`を使用してデフォルト値をマージ。ゼロ値のフィールドにはデフォルト値が適用される。

**テンプレート関数:**
- sprig関数（`Masterminds/sprig`）をすべて利用可能
- `url` - 文字列を`template.URL`に変換
- `safe` - 文字列を`template.HTML`に変換（コメント保持用）

**RTL対応:**
`Hermes.TextDirection`を`TDRightToLeft`に設定することで右から左へのテキスト方向に対応。

### サンプル（examples/）

`examples/`フォルダには以下のサンプルが含まれる:
- `welcome.go` - ウェルカムメール（ボタン付き）
- `invite_code.go` - 招待コード表示
- `receipt.go` - レシート（テーブル使用）
- `reset.go` - パスワードリセット
- `maintenance.go` - メンテナンス通知

`main.go`が全サンプルを各テーマで生成し、`default/`と`flat/`フォルダにHTML/TXTファイルを出力。環境変数を設定すれば実際のメール送信も可能。
