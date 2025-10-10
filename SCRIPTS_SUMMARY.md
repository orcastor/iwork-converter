# 脚本总结

我已经为你创建了以下脚本来帮助管理proto文件生成和codegen：

## 主要脚本

### 1. `generate_all.sh` - 完整构建脚本
**功能：** 一键生成所有proto文件和codegen文件
**使用：** `./generate_all.sh`
**包含：**
- 安装必要工具
- 清理旧文件
- 生成所有proto文件
- 生成所有codegen文件
- 整理导入
- 编译项目

### 2. `generate_proto.sh` - Proto文件生成脚本
**功能：** 仅生成proto文件
**使用：** `./generate_proto.sh`
**包含：**
- 安装protoc-gen-go
- 清理旧pb.go文件
- 生成所有proto文件

### 3. `regenerate_codegen.sh` - Codegen重新生成脚本
**功能：** 重新生成codegen文件（修改JSON配置后使用）
**使用：** `./regenerate_codegen.sh`
**包含：**
- 安装goimports
- 生成所有codegen文件
- 整理导入
- 编译项目

### 4. `quick_test.sh` - 快速测试脚本
**功能：** 快速测试转换功能
**使用：** `./quick_test.sh <input_file>`
**示例：** `./quick_test.sh a.key`

## 辅助文件

### 5. `BUILD_SCRIPTS.md` - 详细使用说明
包含所有脚本的详细使用说明和故障排除指南。

## 使用流程

### 首次设置
```bash
./generate_all.sh
```

### 日常开发
```bash
# 修改了JSON配置后
./regenerate_codegen.sh

# 修改了proto文件后
./generate_proto.sh
./regenerate_codegen.sh

# 快速测试
./quick_test.sh a.key
```

## 脚本特点

1. **自动化：** 自动安装依赖工具
2. **错误处理：** 使用`set -e`确保错误时停止
3. **路径处理：** 自动处理goimports路径问题
4. **清理功能：** 自动清理旧文件
5. **详细输出：** 提供清晰的进度信息

## 注意事项

- 确保已安装`protoc`工具
- 脚本会自动处理Go工具路径问题
- 所有脚本都有错误处理机制
- 建议在修改配置后使用相应的脚本重新生成
