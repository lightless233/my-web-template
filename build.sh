#!/bin/bash

set -euo pipefail

echo "[***] Start building my-web-template ..."
echo "[***] Start collect versions ..."
builtAt="$(date +'%F %T %z')"
builtAtTimestamp="$(date +'%s')"
goVersion=$(go version | sed 's/go version //')
gitAuthor=$(git show -s --format='format:%aN <%ae>' HEAD)
gitCommit=$(git rev-parse HEAD)

ldflags="\
-X 'my-web-template/internal/version.BuiltAt=$builtAt' \
-X 'my-web-template/internal/version.BuiltAtTimestamp=$builtAtTimestamp' \
-X 'my-web-template/internal/version.GoVersion=$goVersion' \
-X 'my-web-template/internal/version.GitAuthor=$gitAuthor' \
-X 'my-web-template/internal/version.GitCommit=$gitCommit' \
-X 'my-web-template/internal/version.AppVersion=$(tr -d "[:space:]" < ./version.txt)' \
"
echo "[***] Versions collected. ldflags=$ldflags"

echo "[***] Start building binary ..."
mkdir -p ./_build/
# 获取操作系统（转换为小写，兼容 macOS/Linux）
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"

# 获取架构并转换为 Go 的格式
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  GOARCH="amd64" ;;
  aarch64) GOARCH="arm64" ;;
  arm64)   GOARCH="arm64" ;;  # Apple Silicon
  *)       echo "Unsupported架构: $ARCH"; exit 1 ;;
esac

# 根据操作系统设置 GOOS
case "$OS" in
  linux*)  GOOS="linux" ;;
  darwin*) GOOS="darwin" ;;
  windows*) GOOS="windows" ;;
  *)       echo "Unsupported系统: $OS"; exit 1 ;;
esac
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
  -ldflags "$ldflags" \
  -o "./_build/my-web-template" cmd/web/web.go
echo "[***] go build done."

# 根据 env_type 环境变量，复制不同的配置文件过去
echo "[***] copy config file, current env_type: $ENV_TYPE"
case "$ENV_TYPE" in
  "testing")
    cp ./config.daily.toml ./_build/config.toml
    ;;
  "staging")
    cp ./config.staging.toml ./_build/config.toml
    ;;
  "production")
    cp ./config.prod.toml ./_build/config.toml
    ;;
  *)
    echo "Unsupported env_type: $ENV_TYPE"
    exit 1
    ;;
esac
