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
    --proto_path=proto \
    proto/TSPMessages.proto proto/TSPArchiveMessages.proto proto/TSPDatabaseMessages.proto

# TSK相关
echo "生成TSK相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSKArchives.proto=github.com/orcastor/iwork-converter/proto/TSK \
    --proto_path=proto \
    proto/TSKArchives.proto

# TSS相关
echo "生成TSS相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSSArchives.proto=github.com/orcastor/iwork-converter/proto/TSS \
    --go_opt=Mproto/TSSArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSS \
    --proto_path=proto \
    proto/TSSArchives.proto proto/TSSArchives.sos.proto

# TSD相关
echo "生成TSD相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSDArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    --go_opt=Mproto/TSDArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSD \
    --go_opt=Mproto/TSDCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSD \
    --proto_path=proto \
    proto/TSDArchives.proto proto/TSDArchives.sos.proto proto/TSDCommandArchives.proto

# TSWP相关
echo "生成TSWP相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSWPArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    --go_opt=Mproto/TSWPArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    --go_opt=Mproto/TSWPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSWP \
    --proto_path=proto \
    proto/TSWPArchives.proto proto/TSWPArchives.sos.proto proto/TSWPCommandArchives.proto

# TSCK相关
echo "生成TSCK相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCKArchives.proto=github.com/orcastor/iwork-converter/proto/TSCK \
    --proto_path=proto \
    proto/TSCKArchives.proto

# TSCE相关
echo "生成TSCE相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCEArchives.proto=github.com/orcastor/iwork-converter/proto/TSCE \
    --proto_path=proto \
    proto/TSCEArchives.proto

# TST相关
echo "生成TST相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSTArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TST \
    --go_opt=Mproto/TSTStylePropertyArchiving.proto=github.com/orcastor/iwork-converter/proto/TST \
    --proto_path=proto \
    proto/TSTArchives.proto proto/TSTArchives.sos.proto proto/TSTCommandArchives.proto proto/TSTStylePropertyArchiving.proto

# TSCH相关
echo "生成TSCH相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHArchives.Common.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCHCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --go_opt=Mproto/TSCH3DArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH \
    --proto_path=proto \
    proto/TSCHArchives.proto proto/TSCHArchives.Common.proto proto/TSCHArchives.sos.proto proto/TSCHCommandArchives.proto proto/TSCH3DArchives.proto

# TSCH Generated
echo "生成TSCH Generated文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHArchives.GEN.proto=github.com/orcastor/iwork-converter/proto/TSCH_Generated \
    --proto_path=proto \
    proto/TSCHArchives.GEN.proto

# TSCH PreUFF
echo "生成TSCH PreUFF文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSCHPreUFFArchives.proto=github.com/orcastor/iwork-converter/proto/TSCH/PreUFF \
    --proto_path=proto \
    proto/TSCHPreUFFArchives.proto

# TSA相关
echo "生成TSA相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TSAArchives.proto=github.com/orcastor/iwork-converter/proto/TSA \
    --go_opt=Mproto/TSACommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TSA \
    --proto_path=proto \
    proto/TSAArchives.proto proto/TSACommandArchives.sos.proto

# KN相关
echo "生成KN相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/KNArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/KN \
    --go_opt=Mproto/KNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/KN \
    --proto_path=proto \
    proto/KNArchives.proto proto/KNArchives.sos.proto proto/KNCommandArchives.proto proto/KNCommandArchives.sos.proto

# TN相关
echo "生成TN相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TNArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TN \
    --go_opt=Mproto/TNCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TN \
    --proto_path=proto \
    proto/TNArchives.proto proto/TNArchives.sos.proto proto/TNCommandArchives.proto proto/TNCommandArchives.sos.proto

# TP相关
echo "生成TP相关文件..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go_opt=Mproto/TPArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    --go_opt=Mproto/TPCommandArchives.proto=github.com/orcastor/iwork-converter/proto/TP \
    --go_opt=Mproto/TPCommandArchives.sos.proto=github.com/orcastor/iwork-converter/proto/TP \
    --proto_path=proto \
    proto/TPArchives.proto proto/TPCommandArchives.proto proto/TPCommandArchives.sos.proto

# 移动生成的文件到正确位置
echo "移动生成的文件到正确位置..."

# 移动TSP文件
if [ -f "TSPMessages.pb.go" ]; then
    echo "移动TSP文件..."
    mkdir -p proto/TSP
    mv TSPMessages.pb.go proto/TSP/
fi
if [ -f "TSPArchiveMessages.pb.go" ]; then
    mv TSPArchiveMessages.pb.go proto/TSP/
fi
if [ -f "TSPDatabaseMessages.pb.go" ]; then
    mv TSPDatabaseMessages.pb.go proto/TSP/
fi

# 移动TSK文件
if [ -f "TSKArchives.pb.go" ]; then
    echo "移动TSK文件..."
    mkdir -p proto/TSK
    mv TSKArchives.pb.go proto/TSK/
fi

# 移动TSS文件
if [ -f "TSSArchives.pb.go" ]; then
    echo "移动TSS文件..."
    mkdir -p proto/TSS
    mv TSSArchives.pb.go proto/TSS/
fi
if [ -f "TSSArchives.sos.pb.go" ]; then
    mv TSSArchives.sos.pb.go proto/TSS/
fi

# 移动TSD文件
if [ -f "TSDArchives.pb.go" ]; then
    echo "移动TSD文件..."
    mkdir -p proto/TSD
    mv TSDArchives.pb.go proto/TSD/
fi
if [ -f "TSDArchives.sos.pb.go" ]; then
    mv TSDArchives.sos.pb.go proto/TSD/
fi
if [ -f "TSDCommandArchives.pb.go" ]; then
    mv TSDCommandArchives.pb.go proto/TSD/
fi

# 移动TSWP文件
if [ -f "TSWPArchives.pb.go" ]; then
    echo "移动TSWP文件..."
    mkdir -p proto/TSWP
    mv TSWPArchives.pb.go proto/TSWP/
fi
if [ -f "TSWPArchives.sos.pb.go" ]; then
    mv TSWPArchives.sos.pb.go proto/TSWP/
fi
if [ -f "TSWPCommandArchives.pb.go" ]; then
    mv TSWPCommandArchives.pb.go proto/TSWP/
fi

# 移动TSCK文件
if [ -f "TSCKArchives.pb.go" ]; then
    echo "移动TSCK文件..."
    mkdir -p proto/TSCK
    mv TSCKArchives.pb.go proto/TSCK/
fi

# 移动TSCE文件
if [ -f "TSCEArchives.pb.go" ]; then
    echo "移动TSCE文件..."
    mkdir -p proto/TSCE
    mv TSCEArchives.pb.go proto/TSCE/
fi

# 移动TST文件
if [ -f "TSTArchives.pb.go" ]; then
    echo "移动TST文件..."
    mkdir -p proto/TST
    mv TSTArchives.pb.go proto/TST/
fi
if [ -f "TSTArchives.sos.pb.go" ]; then
    mv TSTArchives.sos.pb.go proto/TST/
fi
if [ -f "TSTCommandArchives.pb.go" ]; then
    mv TSTCommandArchives.pb.go proto/TST/
fi
if [ -f "TSTStylePropertyArchiving.pb.go" ]; then
    mv TSTStylePropertyArchiving.pb.go proto/TST/
fi

# 移动TSCH文件
if [ -f "TSCHArchives.pb.go" ]; then
    echo "移动TSCH文件..."
    mkdir -p proto/TSCH
    mv TSCHArchives.pb.go proto/TSCH/
fi
if [ -f "TSCHArchives.Common.pb.go" ]; then
    mv TSCHArchives.Common.pb.go proto/TSCH/
fi
if [ -f "TSCHArchives.sos.pb.go" ]; then
    mv TSCHArchives.sos.pb.go proto/TSCH/
fi
if [ -f "TSCHCommandArchives.pb.go" ]; then
    mv TSCHCommandArchives.pb.go proto/TSCH/
fi
if [ -f "TSCH3DArchives.pb.go" ]; then
    mv TSCH3DArchives.pb.go proto/TSCH/
fi

# 移动TSCH Generated文件
if [ -f "TSCHArchives.GEN.pb.go" ]; then
    echo "移动TSCH Generated文件..."
    mkdir -p proto/TSCH_Generated
    mv TSCHArchives.GEN.pb.go proto/TSCH_Generated/
fi

# 移动TSCH PreUFF文件
if [ -f "TSCHPreUFFArchives.pb.go" ]; then
    echo "移动TSCH PreUFF文件..."
    mkdir -p proto/TSCH/PreUFF
    mv TSCHPreUFFArchives.pb.go proto/TSCH/PreUFF/
fi

# 移动TSA文件
if [ -f "TSAArchives.pb.go" ]; then
    echo "移动TSA文件..."
    mkdir -p proto/TSA
    mv TSAArchives.pb.go proto/TSA/
fi
if [ -f "TSACommandArchives.sos.pb.go" ]; then
    mv TSACommandArchives.sos.pb.go proto/TSA/
fi

# 移动KN文件
if [ -f "KNArchives.pb.go" ]; then
    echo "移动KN文件..."
    mkdir -p proto/KN
    mv KNArchives.pb.go proto/KN/
fi
if [ -f "KNArchives.sos.pb.go" ]; then
    mv KNArchives.sos.pb.go proto/KN/
fi
if [ -f "KNCommandArchives.pb.go" ]; then
    mv KNCommandArchives.pb.go proto/KN/
fi
if [ -f "KNCommandArchives.sos.pb.go" ]; then
    mv KNCommandArchives.sos.pb.go proto/KN/
fi

# 移动TN文件
if [ -f "TNArchives.pb.go" ]; then
    echo "移动TN文件..."
    mkdir -p proto/TN
    mv TNArchives.pb.go proto/TN/
fi
if [ -f "TNArchives.sos.pb.go" ]; then
    mv TNArchives.sos.pb.go proto/TN/
fi
if [ -f "TNCommandArchives.pb.go" ]; then
    mv TNCommandArchives.pb.go proto/TN/
fi
if [ -f "TNCommandArchives.sos.pb.go" ]; then
    mv TNCommandArchives.sos.pb.go proto/TN/
fi

# 移动TP文件
if [ -f "TPArchives.pb.go" ]; then
    echo "移动TP文件..."
    mkdir -p proto/TP
    mv TPArchives.pb.go proto/TP/
fi
if [ -f "TPCommandArchives.pb.go" ]; then
    mv TPCommandArchives.pb.go proto/TP/
fi
if [ -f "TPCommandArchives.sos.pb.go" ]; then
    mv TPCommandArchives.sos.pb.go proto/TP/
fi

echo "=== Proto文件生成完成 ==="
