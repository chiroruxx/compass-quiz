# COMPASS 課題

MySQL で生徒・クラス・教師の情報を管理し、Go と Gin で生徒一覧 API を実装する。

## 開発環境

Docker Compose で MySQL と開発用 API を起動する。Go ファイルを保存すると Air が API を再ビルドして再起動する。

```sh
docker compose up -d
```

### テスト

API コンテナでユニットテストと HTTP テストをまとめて実行する。

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

JVM 系言語の経験がまだ少ないため、Go を使用する。OpenAPI から API の型と Gin のルーティングを生成する。構成と各レイヤーの責務は [アーキテクチャ設計](ARCHITECTURE.md)を参照する。
