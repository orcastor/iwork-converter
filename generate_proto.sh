#!/bin/bash

# 生成proto文件的脚本
# 使用方法: ./generate_proto.sh

set -e

echo "=== 生成proto文件 ==="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# 确保protoc-gen-go已安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# 清理旧的pb.go文件（保留TSP相关文件）
echo "清理旧的pb.go文件..."
find proto -name "*.pb.go" -not -path "*/TSP/*" -delete

# 生成基础proto文件
echo "生成基础proto文件..."

# TSP相关
echo "生成TSP相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSPMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    --go_opt=Mproto/TSPArchiveMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    --go_opt=Mproto/TSPDatabaseMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    proto/TSPMessages.proto proto/TSPArchiveMessages.proto proto/TSPDatabaseMessages.proto

# TSK相关
echo "生成TSK相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSKArchives.proto=github.com/orcastor/iwork-converter/proto/TSK \
    proto/TSKArchives.proto

# TSS相关
echo "生成TSS相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSSArchives.proto=github.com/orcastor/iwork-converter/proto/TSS \
    --go_opt=Mproto/TSSArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSS \
    proto/TSSArchives.proto proto/TSSArchives.sos.proto

# TSD相关
echo "生成TSD相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSDArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    --go_opt=Mproto/TSDArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSD \
    --go_opt=Mproto/TSDCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    proto/TSDArchives.proto proto/TSDArchives.sos.proto proto/TSDCommandArchives.proto

# TSWP相关
echo "生成TSWP相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSWPArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    --go_opt=Mproto/TSWPArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    --go_opt=Mproto/TSWPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    proto/TSWPArchives.proto proto/TSWPArchives.sos.proto proto/TSWPCommandArchives.proto

# TSCK相关
echo "生成TSCK相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCKArchives.proto=github.com/orcastor/iwork-converter/proto/TSCK \
    proto/TSCKArchives.proto

# TSCE相关
echo "生成TSCE相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCEArchives.proto=github.com/orcastor/iwork-converter/proto/TSCE \
    proto/TSCEArchives.proto

# TST相关
echo "生成TST相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSTArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTStylePropertyArchiving.proto=github.com/orcastor/iwork-converter/proto/TST \
    proto/TSTArchives.proto proto/TSTArchives.sos.proto proto/TSTCommandArchives.proto proto/TSTStylePropertyArchiving.proto

# TSCH相关
echo "生成TSCH相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHArchives.Common.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCH3DArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    proto/TSCHArchives.proto proto/TSCHArchives.Common.proto proto/TSCHArchives.sos.proto proto/TSCHCommandArchives.proto proto/TSCH3DArchives.proto

# TSCH Generated
echo "生成TSCH Generated文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHArchives.GEN.proto=github.com/orcastor/iwork-converter/proto/TSCH_Generated \
    proto/TSCHArchives.GEN.proto

# TSCH PreUFF
echo "生成TSCH PreUFF文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHPreUFFArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH/PreUFF \
    proto/TSCHPreUFFArchives.proto

# TSA相关
echo "生成TSA相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSAArchives.proto=github.com/orcastor/iwork-converter/proto/TSA \
    --go_opt=Mproto/TSACommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSA \
    proto/TSAArchives.proto proto/TSACommandArchives.sos.proto

# KN相关
echo "生成KN相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/KNArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    proto/KNArchives.proto proto/KNArchives.sos.proto proto/KNCommandArchives.proto proto/KNCommandArchives.sos.proto

# TN相关
echo "生成TN相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TNArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    proto/TNArchives.proto proto/TNArchives.sos.proto proto/TNCommandArchives.proto proto/TNCommandArchives.sos.proto

# TP相关
echo "生成TP相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TPArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    --go_opt=Mproto/TPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    --go_opt=Mproto/TPCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TP \
    proto/TPArchives.proto proto/TPCommandArchives.proto proto/TPCommandArchives.sos.proto

echo "=== Proto文件生成完成 ==="
