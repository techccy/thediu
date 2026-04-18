#!/bin/bash

# 创建一个会报错的测试文件

cat > /tmp/test_error.py << 'EOF'
import os

config_file = open("config.json", "r")
content = config_file.read()
print(content)
config_file.close()
EOF

echo "=== 测试深度上下文感知 ==="
echo ""
echo "1. 创建测试文件 /tmp/test_error.py"
cat /tmp/test_error.py
echo ""

echo "2. 执行测试脚本（预期会报错）"
cd /tmp
python3 test_error.py 2>&1 | head -5
echo ""

echo "3. 检查目录结构"
tree -L 1 /tmp 2>/dev/null || ls -la /tmp | head -10
echo ""

echo "=== 测试完成 ==="
