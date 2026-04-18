# CCY Assistant

受 [thefuck](https://github.com/nvbn/thefuck) 启发，一个 AI 驱动的命令行助手。当你在终端遇到错误时，它会自动分析错误并提供解决方案。

## 功能特性

- 🔍 **智能错误分析** - 利用 AI 深度分析命令错误
- 🎯 **多 AI 提供商支持** - 集成 OpenAI、DeepSeek、Ollama
- 🚀 **本地缓存记忆** - 自动缓存成功的修复方案，下次直接使用
- 📁 **上下文感知** - 自动收集目录结构和相关文件内容
- 🎨 **美观的 TUI 界面** - 基于 Bubbletea 的现代化终端界面
- ✅ **安全确认机制** - 执行前需要用户确认，支持黑名单保护
- 🔄 **自动故障转移** - 网络失败自动切换到本地 Ollama
- 🛡️ **隐私保护** - 自动过滤敏感文件（.env, .key, 凭证等）
- 🌍 **多平台支持** - 支持 macOS, Linux

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
sudo mv ccy-core /usr/local/bin/

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

### 配置 AI 提供商

#### 配置文件

CCY Assistant 会自动在 `~/.ccy/config.yaml` 创建配置文件：

```yaml
default_provider: ollama  # 默认提供商
providers:
  openai:
    base_url: https://api.openai.com/v1
    api_key: env:CCY_OPENAI_KEY
    model: gpt-4o-mini
  deepseek:
    base_url: https://api.deepseek.com
    api_key: env:CCY_DEEPSEEK_KEY
    model: deepseek-chat
  ollama:
    base_url: http://localhost:11434
    api_key: ""
    model: qwen2.5:7b
```

#### 设置 API Key

使用环境变量：

```bash
# OpenAI
export CCY_OPENAI_KEY=sk-your-key-here

# DeepSeek
export CCY_DEEPSEEK_KEY=sk-your-key-here
```

或使用 `ccy-core --config` 查看详细配置指南

#### 切换提供商

```bash
ccy-core --switch
```

### 命令行选项

```bash
ccy-core --init      # 生成 Shell 初始化脚本
ccy-core --config    # 显示配置指南
ccy-core --switch    # 切换 AI 提供商
ccy-core --help      # 显示帮助信息
```

## 工作原理

1. **错误捕获** - Shell 函数捕获上一条失败的命令和错误信息
2. **上下文收集** - 自动收集目录结构，分析错误信息中的文件路径并提取相关内容
3. **AI 分析** - 发送到 AI API 分析并生成修复方案
4. **智能缓存** - 检查本地缓存，如有成功记录直接返回
5. **安全确认** - 展示分析结果和修复命令，等待用户确认
6. **命令执行** - 用户确认后执行修复命令并记录结果

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
├── internal/
│   ├── config/          # 配置管理
│   ├── context/         # 上下文收集
│   ├── memory/          # 本地缓存 (SQLite)
│   ├── model/           # 数据模型
│   ├── provider/        # AI 提供商 (OpenAI, DeepSeek, Ollama)
│   ├── shellinit/       # Shell 初始化脚本
│   └── tui/             # 终端 UI (Bubbletea)
├── PRD/                 # 产品需求文档
├── .github/             # GitHub Actions
│   └── workflows/       # CI/CD 配置
└── .goreleaser.yaml     # GoReleaser 配置
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

## 安全特性

- **命令黑名单** - 自动保护 rm、mkfs、reboot 等危险命令
- **敏感文件过滤** - 不读取 .env、.key、凭证等敏感文件
- **用户确认** - 所有命令执行前都需要用户确认
- **超时机制** - 防止命令执行卡住

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License