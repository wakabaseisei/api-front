# api-front

本リポジトリは、MSプロジェクトにおける **BFF（Backend For Frontend）サービス** を提供するアプリケーションです。  
フロントエンドに最適化されたAPIを集約的に提供し、複数のバックエンドサービスとの橋渡しを担います。

---

## 🧩 概要

- **通信方式:** Connect RPC を採用し、Protobuf による型安全・高速な通信を実現
- **BFF設計:** フロントエンドに最適化したAPIインターフェースを提供
- **CORS対応:** `connectrpc.com/cors` + `rs/cors` によるクロスオリジン通信制御に対応

---

## 📁 ディレクトリ構成

```
.
├── .github/workflows/ # GitHub Actions ワークフロー（ECRビルド＆プッシュ）
├── internal/
│ ├── cmd/app/ # アプリケーションのエントリーポイント（main.go）
│ ├── config/ # 環境変数パース処理（caarlos0/env）
│ ├── domain/ # ドメイン層（モデル、サービスインターフェース）
│ ├── driver/
│ │ ├── client/ # 外部サービス（ms-user）との連携クライアント
│ │ └── grpc/ # Connect RPCハンドラの実装
│ └── usecase/ # ユースケース層（アプリケーションサービス）
├── go.mod / go.sum # Go モジュール管理ファイル
├── Dockerfile # アプリケーションのビルド＆実行イメージ（distroless使用）
└── README.md
```

---

## 🚀 ビルド・デプロイ

### GitHub Actions による CI/CD

- `main` ブランチに Push されると、以下を自動実行：
  1. Docker イメージをビルド（BuildKit 使用）
  2. ECR に `dev-日付-コミットSHA` 形式でタグ付き Push
  3. キャッシュ利用による高速なビルド

---

## 🌐 API仕様

- Connect RPC + Protobuf によってスキーマ定義
- スキーマファイルは [`ms-protobuf`](https://github.com/wakabaseisei/ms-protobuf) にて一元管理
- 対応メソッド例：
  - `Greet(name: string): string`
  - `Ping(): void`

---

## 📄 備考

- `USER_SERVICE_ENDPOINT` を環境変数で指定し、`ms-user` との接続先を柔軟に切替可能  
- Aurora への直接アクセスは行わず、永続化処理は `ms-user` をはじめとする各マイクロサービスに委譲する構成です。これにより、データアクセスの責務をドメインごとのサービスに分離し、今後の機能追加やスケーリングにも柔軟に対応できます。
