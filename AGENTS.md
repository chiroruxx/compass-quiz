# Compass Quiz の作業ガイド

## 参照資料

- [課題2の要件](docs/tasks/2_api/requirements.md)には、生徒一覧 API と提出物の要件がある。
- [課題2の実装ステップ](docs/tasks/2_api/steps.md)には、実装順序とステップごとの完了条件がある。
- [README](README.md)と[アーキテクチャ設計](ARCHITECTURE.md)には、API の実装方針と開発・実行環境の設計がある。
- [課題1の要件](docs/tasks/1_db/requirements.md)と[DB 設計](docs/design/db.md)には、API が利用する既存の MySQL の要件と構造がある。

## 作業の進め方

1. 着手前に、依頼されたステップの作業項目と完了条件を読み、関連する要件と設計を確認する。判断に必要な仕様が資料から決まらない場合は、ユーザーに確認する。
2. 依頼されたステップに必要な変更を行う。動作確認に必要なテストケースは、その機能を実装するステップで追加する。
3. ステップごとの完了条件を項目ごとに確認し、実行した確認と未確認の項目を報告する。起動・テスト・コード生成のコマンドは、リポジトリで実行できるようになった段階でここに追記する。

## レビュー・PR・CI

1. 各ステップの完了条件を確認したら、コード変更には [実装ステップのレビュースキル `$review-step`](.agents/skills/review-step/SKILL.md)、Markdown の変更には [ドキュメントのレビュースキル `$review-markdown`](.agents/skills/review-markdown/SKILL.md) を使う。スキルの指摘のうち修正できるものは修正し、人間の判断が必要なものは報告する。すべての変更について人間のレビューを受け、必要なレビューを通過したことを確認する。
2. レビュー通過後に PR を作成し、CI を実行して結果を確認する。PR は人間が差分を追いやすい規模にする。複数の PR が必要な場合は、先行 PR のブランチを後続 PR のベースにして作成する。[Stacked Pull Requests](https://docs.github.com/en/pull-requests/get-started/about-stacked-prs)

## マージ後の改善

- 実装 PR のマージ後に、作業中に得た知見やレビューでの指摘を振り返り、AI の設定、ハーネス、ガードレールを改善できる箇所を探す。複数の作業に共通して適用できる、システム全体の改善を対象とする。
- 改善箇所が見つかった場合は、別途改善 PR を作成する。
- 振り返りに役立つ気づきは、作業中に `docs/2_api/worklog.md` に記録してよい。

## Go と生成コード

- Go コードを書く前に、[`go.mod`](go.mod) の Go バージョンに対応する [JetBrains Modern Go Guidelines](https://github.com/JetBrains/go-modern-guidelines) を確認する。Codex では `modern-go-guidelines` プラグインの `use-modern-go` スキルを使用する。
- 生成コードを変更するときは、生成元を編集してから再生成する。API の型と Gin のルーティングは OpenAPI 定義から `oapi-codegen` で、データアクセスコードは MySQL のスキーマと SQL クエリから `sqlc` で、テスト用モックは対象のインターフェースから gomock で生成する。
