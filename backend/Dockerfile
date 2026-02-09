# Goの軽量イメージを使用
FROM golang:1.25-alpine

# コンテナ内の作業ディレクトリ
WORKDIR /app

# ライブリロード（コード変更時に自動再起動）ツールをインストール
RUN go install github.com/air-verse/air@latest

# コンテナ内でマイグレーションコマンドを使えるようにツールをインストール
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ソースコードをコピー
COPY . .

# 依存関係の解決（後でgo mod initした後に機能）
RUN go mod download || true

# airコマンドでサーバー起動
CMD ["air", "-c", ".air.toml"]
