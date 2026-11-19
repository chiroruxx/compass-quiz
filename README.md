# COMPASS 課題

## 開発環境

Docker Compose で MySQL と開発用 API を起動する。Go ファイルを保存すると Air が API を再ビルドして再起動する。

```sh
docker compose up -d
```

### テスト

MySQL の起動後、API コンテナでユニットテストと HTTP テストをまとめて実行する。HTTP テストは、MySQL の初回起動時に作成される `compass_quiz_test` を使い、テスト内で用意したデータをロールバックする。

```sh
docker compose exec api go test -count=1 ./...
```

### コード生成と CI

`openapi.yaml` を変更したら、Go の型と Gin のルーティングを再生成する。

```sh
docker compose exec api go generate ./...
```

CI と同じリントを API コンテナで実行できる。

```sh
docker compose exec api go tool -modfile=tools/golangci-lint/go.mod golangci-lint run
docker compose exec api go tool -modfile=tools/actionlint/go.mod actionlint
```

## 実装方針

### 言語と実行環境

- JVM 系言語の経験がまだ少ないため、Go を使用する。
- マルチステージビルドを使い、`golang` イメージで API を静的ビルドし、`scratch` イメージで実行する。通常起動には `compose.yaml` を使う。
- ローカル開発時は `compose.dev.yaml` を明示的に重ねる。API のソースをマウントし、Go ファイルの変更時に Air で再ビルド・再起動する。MySQL と HTTPS の設定は通常起動と共通にする。

### アーキテクチャ

[アーキテクチャ設計](ARCHITECTURE.md)を参照する。

### フィルター、ソート、レスポンス

- 検索条件が複数指定された場合は AND 条件にする。
- 生徒名でソートする場合、同名の生徒は生徒 ID の昇順で並べる。
- 該当件数が 0 件なら、本文なしの `404` を返す。
- 指定ページだけが範囲外なら、`200` を返す。
- 不正なクエリには、本文なしの `400` を返す。

### TLS

- ローカル用の自己署名証明書を初回起動時に生成し、Go の HTTP サーバーで TLS を終端する。

### テスト

以下のテストを用意し、自動で実行できるようにする。

- ユニットテストは DB を使わず、必要に応じてモックを使う。
- HTTP テストは `httptest` とテスト用の MySQL を使い、API 単位で実行する。
