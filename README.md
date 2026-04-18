# DIU Assistant

受 [thefuck](https://github.com/nvbn/thefuck) 启发，一个 AI 驱动的命令行助手。当你在终端遇到错误时，它会自动分析错误并提供解决方案。

## 功能特性

- **智能错误分析** - 利用 AI 深度分析命令错误
- **多 AI 提供商支持** - 集成 OpenAI、DeepSeek、Ollama
- **上下文感知** - 自动收集目录结构和相关文件内容
- **美观的 TUI 界面** - 基于 Bubbletea 的现代化终端界面
- **安全确认机制** - 执行前需要用户确认，支持黑名单保护
- **自动故障转移** - 网络失败自动切换到本地 Ollama
- **隐私保护** - 自动过滤敏感文件（.env, .key, 凭证等）
- **多平台支持** - 支持 macOS, Linux, Windows (x86_64, ARM64)

## 快速开始

### 使用 Homebrew 安装（推荐）

```bash
# 添加 Tap 仓库
brew tap techccy/tap

# 安装 DIU
brew install diu

# 初始化 Shell 集成
diu-core --init >> ~/.zshrc  # zsh
diu-core --init >> ~/.bashrc # bash

# 重新加载 shell 配置
source ~/.zshrc  # zsh
source ~/.bashrc # bash
```

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/techccy/thediu.git
cd thediu
```

```
* 编译 (生成当前平台的二进制文件)
go build -o diu-core
```

* 将以下变量放在环境变量末尾
```
export PATH="{PATH_TO}thediu/build:$PATH"
eval "$(diu-core --init)"
```

* 重载配置
```
source ~/.zshrc
```

```
# Windows (使用 PowerShell)
Move-Item diu-core "$env:USERPROFILE\bin\"
```

```
# 初始化
eval "$(diu-core --init)" >> ~/.zshrc  # zsh
eval "$(diu-core --init)" >> ~/.bashrc # bash
# 重新加载 shell 配置
source ~/.zshrc  # zsh
source ~/.bashrc # bash
```

### 交叉编译

项目支持编译多平台的独立可执行文件：

```bash
# macOS x86_64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o diu-core-darwin-amd64

# macOS ARM64 (Apple Silicon)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o diu-core-darwin-arm64

# Linux x86_64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o diu-core-linux-amd64

# Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o diu-core-linux-arm64

# Windows x86_64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o diu-core-windows-amd64.exe
```

所有二进制文件均为静态链接，可独立运行，无需额外依赖。

### 从 Release 安装

1. 访问 [GitHub Release 页面](https://github.com/techccy/diu-assistant/releases)
2. 根据你的系统下载对应的二进制文件：

| 平台 | 架构 | 文件名 |
|------|------|--------|
| macOS | x86_64 (Intel) | `diu-core_darwin_amd64.tar.gz` |
| macOS | ARM64 (Apple Silicon) | `diu-core_darwin_arm64.tar.gz` |
| Linux | x86_64 | `diu-core_linux_amd64.tar.gz` |
| Linux | ARM64 | `diu-core_linux_arm64.tar.gz` |
| Windows | x86_64 | `diu-core_windows_amd64.zip` |

3. 解压并移动到 PATH 目录
4. 运行初始化命令：`diu-core --init`

## 使用方法

### 基本使用

当你在终端遇到错误时：

```bash
# 假设你运行了一个失败的命令
git push origin main
# fatal: The current branch has no upstream branch

# 运行 diu 来获取解决方案
diu
```

### Shell 集成

运行以下命令将 diu 集成到你的 shell 中：

```bash
eval "$(diu-core --init)"
```

这会将一个 `diu` 函数添加到你的 shell 中，它会：
1. 自动捕获上一条失败的命令
2. 分析错误信息
3. 提供修复建议
4. 让你确认是否执行

### 配置 AI 提供商

#### 配置文件

DIU Assistant 会自动在 `~/.diu/config.yaml` 创建配置文件：

```yaml
default_provider: ollama  # 默认提供商
providers:
  openai:
    base_url: https://api.openai.com/v1
    api_key: env:DIU_OPENAI_KEY
    model: gpt-4o-mini
  deepseek:
    base_url: https://api.deepseek.com
    api_key: env:DIU_DEEPSEEK_KEY
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
export DIU_OPENAI_KEY=sk-your-key-here

# DeepSeek
export DIU_DEEPSEEK_KEY=sk-your-key-here
```

或使用 `diu-core --config` 查看详细配置指南

#### 切换提供商

```bash
diu-core --switch
```

### 命令行选项

```bash
diu-core --init      # 生成 Shell 初始化脚本
diu-core --config    # 显示配置指南
diu-core --switch    # 切换 AI 提供商
diu-core --help      # 显示帮助信息
```

## 工作原理

1. **错误捕获** - Shell 函数捕获上一条失败的命令和错误信息
2. **上下文收集** - 自动收集目录结构，分析错误信息中的文件路径并提取相关内容
3. **AI 分析** - 发送到 AI API 分析并生成修复方案
4. **安全确认** - 展示分析结果和修复命令，等待用户确认
5. **命令执行** - 用户确认后执行修复命令并记录结果

## 开发路线图

- [x] **Phase 1**: 核心逻辑 MVP (API 交互与 JSON 解析)
- [x] **Phase 2**: Shell 劫持与上下文捕获
- [x] **Phase 3**: 优雅的终端交互界面 (TUI)
- [x] **Phase 4**: CI/CD 自动化构建与多平台支持
- [x] **Phase 5**: Homebrew 一键安装

## 项目结构

```
diu-assistant/
├── main.go              # 主程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── context/         # 上下文收集
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
# 构建当前平台
go build -o diu-core

# 构建所有平台 (使用 GoReleaser)
goreleaser build --snapshot --clean

# 构建并打包 (用于发布)
goreleaser release --clean
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

## 许可证

MIT License