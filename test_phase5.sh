#!/bin/bash

# 测试脚本用于验证 Phase 5 功能

echo "=== Phase 5 功能测试 ==="
echo ""

echo "1. 检查配置文件..."
if [ -f ~/.ccy/config.yaml ]; then
    echo "✅ 配置文件存在"
    echo "内容:"
    cat ~/.ccy/config.yaml
else
    echo "❌ 配置文件不存在"
fi
echo ""

echo "2. 检查历史记录数据库..."
if [ -f ~/.ccy/history.db ]; then
    echo "✅ 历史记录数据库存在"
else
    echo "❌ 历史记录数据库不存在（首次运行时正常）"
fi
echo ""

echo "3. 检查目录结构..."
echo "项目结构:"
find . -type f -name "*.go" | head -20
echo ""

echo "4. 测试帮助信息..."
./ccy-core --help
echo ""

echo "=== 测试完成 ==="
