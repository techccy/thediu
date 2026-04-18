# CCY Assistant

一个智能的命令行助手，当你在终端遇到错误时，它会自动分析错误并提供解决方案。

## 功能特性

- 🤖 AI 驱动的错误分析
- 🎨 美观的终端交互界面 (TUI)
- 🚀 一键执行修复命令
- 🔒 安全确认机制
- 📦 多平台支持 (macOS, Linux)

## 快速开始

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/ccy-ai/ccy-assistant.git
cd ccy-assistant

# 安装依赖
go mod download

# 编译
go build -o ccy-core

# 移动到 PATH 目录
mv ccy-core /usr/local/bin/

# 初始化
eval "$(ccy-core --init)" >> ~/.zshrc
source ~/.zshrc
```

### 从 Release 安装

1. 访问 [GitHub Release 页面](https://github.com/ccy-ai/ccy-assistant/releases)
2. 下载对应平台的二进制文件
3. 解压并移动到 PATH 目录
4. 运行初始化命令

## 使用方法

### 基本使用

当你在终端遇到错误时：

```bash
# 假设你运行了一个失败的命令
git push origin main
# fatal: The current branch has no upstream branch

# 运行 ccy 来获取解决方案
ccy
```

### Shell 集成

运行以下命令将 ccy 集成到你的 shell 中：

```bash
eval "$(ccy-core --init)"
```

这会将一个 `ccy` 函数添加到你的 shell 中，它会：
1. 自动捕获上一条失败的命令
2. 分析错误信息
3. 提供修复建议
4. 让你确认是否执行

### 配置 API Key

创建 `.env` 文件：

```bash
CCY_API_KEY=your_api_key_here
```

或设置环境变量：

```bash
export CCY_API_KEY=your_api_key_here
```

## 开发路线图

- [x] **Phase 1**: 核心逻辑 MVP (API 交互与 JSON 解析)
- [x] **Phase 2**: Shell 劫持与上下文捕获
- [x] **Phase 3**: 优雅的终端交互界面 (TUI)
- [x] **Phase 4**: CI/CD 自动化构建
- [ ] **Phase 5**: Homebrew 一键安装

## 项目结构

```
ccy-assistant/
├── main.go              # 主程序入口
├── internal/            # 内部包
│   ├── api/            # API 客户端
│   ├── config/         # 配置管理
│   ├── model/          # 数据模型
│   ├── shellinit/      # Shell 初始化
│   └── tui/            # 终端 UI
├── PRD/                # 产品需求文档
├── .github/            # GitHub Actions
│   └── workflows/      # CI/CD 配置
└── .goreleaser.yaml    # GoReleaser 配置
```

## 开发

### 运行测试

```bash
go test -v ./...
```

### 构建

```bash
go build -o ccy-core
```

### 代码检查

```bash
go vet ./...
golangci-lint run
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
