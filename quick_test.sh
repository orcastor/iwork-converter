#!/bin/bash

# 快速测试脚本
# 使用方法: ./quick_test.sh [input_file]

set -e

echo "=== 快速测试脚本 ==="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# 检查是否提供了输入文件
if [ $# -eq 0 ]; then
    echo "使用方法: $0 <input_file>"
    echo "例如: $0 a.key"
    exit 1
fi

INPUT_FILE="$1"
OUTPUT_FILE="${INPUT_FILE%.*}.html"

echo "输入文件: $INPUT_FILE"
echo "输出文件: $OUTPUT_FILE"

# 检查输入文件是否存在
if [ ! -f "$INPUT_FILE" ]; then
    echo "错误: 输入文件 '$INPUT_FILE' 不存在"
    exit 1
fi

# 编译项目
echo "编译项目..."
go build

# 运行转换
echo "运行转换..."
./iwork-converter "$INPUT_FILE" "$OUTPUT_FILE"

echo "=== 转换完成 ==="
echo "输出文件: $OUTPUT_FILE"
