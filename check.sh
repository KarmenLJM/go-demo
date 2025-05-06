#!/usr/bin/env bash
set -euo pipefail

rm -rf coverdata
mkdir coverdata

GOCOVERDIR=coverdata \
go run -cover -coverpkg=./... -covermode=atomic . &   # ← 后台启动
PID=$!
echo -e "$PID"

# 等待服务就绪（最多 5 秒）
for i in {1..50}; do
  if curl -s http://localhost:8080/hello >/dev/null 2>&1; then break; fi
  sleep 0.1
done

# 抓取覆盖率
echo -e "\n===== live coverage ====="
curl -s http://localhost:8080/debug/pprof/coverage?format=text | head

kill $PID