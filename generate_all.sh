#!/bin/bash

# 生成所有proto文件和codegen文件的脚本
# 使用方法: ./generate_all.sh

set -e

echo "=== 开始生成所有proto文件和codegen文件 ==="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# 确保protoc-gen-go已安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# 确保goimports已安装
if ! command -v goimports &> /dev/null; then
    echo "安装goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# 清理旧的pb.go文件
echo "清理旧的pb.go文件..."
find proto -name "*.pb.go" -delete

# 生成基础proto文件
echo "生成基础proto文件..."

# TSP相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSPMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    -M proto/TSPArchiveMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    -M proto/TSPDatabaseMessages.proto=github.com/orcastor/iwork-converter/proto/TSP \
    proto/TSPMessages.proto proto/TSPArchiveMessages.proto proto/TSPDatabaseMessages.proto

# TSK相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSKArchives.proto=github.com/orcastor/iwork-converter/proto/TSK \
    proto/TSKArchives.proto

# TSS相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSSArchives.proto=github.com/orcastor/iwork-converter/proto/TSS \
    -M proto/TSSArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSS \
    proto/TSSArchives.proto proto/TSSArchives.sos.proto

# TSD相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSDArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    -M proto/TSDArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSD \
    -M proto/TSDCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    proto/TSDArchives.proto proto/TSDArchives.sos.proto proto/TSDCommandArchives.proto

# TSWP相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSWPArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    -M proto/TSWPArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    -M proto/TSWPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    proto/TSWPArchives.proto proto/TSWPArchives.sos.proto proto/TSWPCommandArchives.proto

# TSCK相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSCKArchives.proto=github.com/orcastor/iwork-converter/proto/TSCK \
    proto/TSCKArchives.proto

# TSCE相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSCEArchives.proto=github.com/orcastor/iwork-converter/proto/TSCE \
    proto/TSCEArchives.proto

# TST相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSTArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    -M proto/TSTArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TST \
    -M proto/TSTCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    -M proto/TSTStylePropertyArchiving.proto=github.com/orcastor/iwork-converter/proto/TST \
    proto/TSTArchives.proto proto/TSTArchives.sos.proto proto/TSTCommandArchives.proto proto/TSTStylePropertyArchiving.proto

# TSCH相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSCHArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    -M proto/TSCHArchives.Common.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    -M proto/TSCHArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    -M proto/TSCHCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    -M proto/TSCH3DArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    proto/TSCHArchives.proto proto/TSCHArchives.Common.proto proto/TSCHArchives.sos.proto proto/TSCHCommandArchives.proto proto/TSCH3DArchives.proto

# TSCH Generated
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSCHArchives.GEN.proto=github.com/orcastor/iwork-converter/proto/TSCH_Generated \
    proto/TSCHArchives.GEN.proto

# TSCH PreUFF
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSCHPreUFFArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH/PreUFF \
    proto/TSCHPreUFFArchives.proto

# TSA相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TSAArchives.proto=github.com/orcastor/iwork-converter/proto/TSA \
    -M proto/TSACommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSA \
    proto/TSAArchives.proto proto/TSACommandArchives.sos.proto

# KN相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/KNArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    -M proto/KNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    -M proto/KNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    -M proto/KNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    proto/KNArchives.proto proto/KNArchives.sos.proto proto/KNCommandArchives.proto proto/KNCommandArchives.sos.proto

# TN相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TNArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    -M proto/TNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    -M proto/TNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    -M proto/TNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    proto/TNArchives.proto proto/TNArchives.sos.proto proto/TNCommandArchives.proto proto/TNCommandArchives.sos.proto

# TP相关
protoc --go_out=. --go_opt=paths=source_relative \
    -M proto/TPArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    -M proto/TPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    -M proto/TPCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TP \
    proto/TPArchives.proto proto/TPCommandArchives.proto proto/TPCommandArchives.sos.proto

echo "Proto文件生成完成"

# 生成codegen文件
echo "生成codegen文件..."

# 生成Common.json的index文件
go run codegen/codegen.go codegen/Common.json Common > index/common.go

# 生成Pages.json的index文件
go run codegen/codegen.go codegen/Pages.json Pages > index/pages.go

# 生成Numbers.json的index文件
go run codegen/codegen.go codegen/Numbers.json Numbers > index/numbers.go

# 生成Keynote.json的index文件
go run codegen/codegen.go codegen/Keynote.json Keynote > index/keynote.go

echo "Codegen文件生成完成"

# 使用goimports整理导入
echo "整理Go导入..."
if command -v goimports &> /dev/null; then
    goimports -w index/*.go
else
    ~/go/bin/goimports -w index/*.go
fi

# 编译项目
echo "编译项目..."
go build

echo "=== 所有文件生成完成 ==="
echo "项目已成功编译"
