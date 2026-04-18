# Phase 4: CI/CD 自动化构建

本阶段已完成 GoReleaser 和 GitHub Actions 配置，实现了自动化构建和发布。

## 配置文件

### 1. GoReleaser 配置 (`.goreleaser.yaml`)

已配置以下功能：
- 支持多平台编译：darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
- 自动生成 SHA256 校验和
- 自动生成 Changelog
- 优化的二进制文件（通过 ldflags 去除调试信息）

### 2. GitHub Actions Workflow (`.github/workflows/release.yml`)

已配置自动触发：
- 当推送以 `v` 开头的标签时（如 `v1.0.0`），自动触发构建
- 自动在 GitHub Release 页面发布所有平台的二进制包
- 支持 GPG 签名（可选）

### 3. CI Workflow (`.github/workflows/ci.yml`)

已配置持续集成：
- 自动运行代码检查 (go vet)
- 自动运行测试
- 跨平台构建验证 (Linux, macOS)
- 代码 Lint 检查 (golangci-lint)

## 使用说明

### 触发发布流程

1. 更新版本号和更新日志（可选）
2. 创建并推送标签：

```bash
git tag v1.0.0
git push origin v1.0.0
```

3. GitHub Actions 会自动触发，编译所有平台的二进制文件
4. 构建完成后，二进制文件会自动发布到 GitHub Release 页面

### 发布的文件

每次发布会生成以下文件：
- `ccy-assistant_v1.0.0_darwin_amd64.tar.gz` - Mac Intel 版本
- `ccy-assistant_v1.0.0_darwin_arm64.tar.gz` - Mac M1/M2/M3 版本
- `ccy-assistant_v1.0.0_linux_amd64.tar.gz` - Linux AMD64 版本
- `ccy-assistant_v1.0.0_linux_arm64.tar.gz` - Linux ARM64 版本
- `ccy-assistant_v1.0.0_SHA256SUMS` - 所有文件的 SHA256 校验和

### 安装二进制文件

用户可以从 GitHub Release 页面下载对应平台的二进制文件：

```bash
# 下载对应的压缩包
wget https://github.com/ccy-ai/ccy-assistant/releases/download/v1.0.0/ccy-assistant_v1.0.0_darwin_arm64.tar.gz

# 解压
tar -xzf ccy-assistant_v1.0.0_darwin_arm64.tar.gz

# 移动到 PATH 目录
mv ccy-core /usr/local/bin/

# 初始化
eval "$(ccy-core --init)" >> ~/.zshrc
source ~/.zshrc
```

### 本地测试构建

在推送标签前，可以先在本地测试构建：

```bash
# 安装 goreleaser
brew install goreleaser

# 测试构建（不发布）
goreleaser build --clean --snapshot

# 完整测试（不发布）
goreleaser release --clean --snapshot --skip-publish
```

## 下一步

Phase 4 已完成。下一步是 **Phase 5：实现 Homebrew 一键安装**，包括：
1. 创建 homebrew-tap 仓库
2. 在 GoReleaser 中配置 Homebrew 模块
3. 用户可以通过 `brew tap ccy-ai/tap && brew install ccy` 一键安装
