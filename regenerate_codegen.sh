#!/bin/bash

# 重新生成codegen文件的脚本
# 使用方法: ./regenerate_codegen.sh

set -e

echo "=== 重新生成codegen文件 ==="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# 确保goimports已安装
if ! command -v goimports &> /dev/null; then
    echo "安装goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# 生成codegen文件
echo "生成codegen文件..."

# 生成Common.json的index文件
echo "生成 index/common.go..."
go run codegen/codegen.go codegen/Common.json Common > index/common.go

# 生成Pages.json的index文件
echo "生成 index/pages.go..."
go run codegen/codegen.go codegen/Pages.json Pages > index/pages.go

# 生成Numbers.json的index文件
echo "生成 index/numbers.go..."
go run codegen/codegen.go codegen/Numbers.json Numbers > index/numbers.go

# 生成Keynote.json的index文件
echo "生成 index/keynote.go..."
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

echo "=== Codegen文件重新生成完成 ==="
echo "项目已成功编译"
