package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	styles *styles
}

type styles struct {
	title    lipgloss.Style
	analysis lipgloss.Style
	command  lipgloss.Style
	box      lipgloss.Style
	error    lipgloss.Style
	success  lipgloss.Style
	spinner  lipgloss.Style
}

func NewUI() *UI {
	styles := newStyles()
	return &UI{styles: styles}
}

func newStyles() *styles {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true).
		MarginBottom(1)

	analysisStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginBottom(1)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("50")).
		Background(lipgloss.Color("236")).
		Bold(true).
		Padding(0, 2).
		MarginTop(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("105")).
		Padding(1, 2).
		Margin(1, 0)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		MarginBottom(1)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	spinnerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		MarginLeft(1)

	return &styles{
		title:    titleStyle,
		analysis: analysisStyle,
		command:  commandStyle,
		box:      boxStyle,
		error:    errorStyle,
		success:  successStyle,
		spinner:  spinnerStyle,
	}
}

func (ui *UI) ShowLoading(message string) func() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	stop := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-stop:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r%s %s%s", frames[i%len(frames)], message, "...")
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		close(stop)
		time.Sleep(50 * time.Millisecond)
		fmt.Print("\r\033[K")
	}
}

func (ui *UI) ShowResult(analysis, command string) {
	title := ui.styles.title.Render("🔍 AI 分析结果")
	analysisText := ui.styles.analysis.Render(analysis)
	commandText := ui.styles.command.Render(fmt.Sprintf("▶ %s", command))

	content := lipgloss.JoinVertical(lipgloss.Left, analysisText, commandText)
	box := ui.styles.box.Render(content)

	fmt.Println()
	fmt.Println(title)
	fmt.Println(box)
}

func (ui *UI) AskConfirmation() (bool, error) {
	var action string
	err := huh.NewSelect[string]().
		Title("请选择操作").
		Options(
			huh.NewOption("🚀 执行", "execute"),
			huh.NewOption("❌ 取消", "cancel"),
		).
		Value(&action).
		Run()
	if err != nil {
		return false, err
	}
	return action == "execute", nil
}

func (ui *UI) ExecuteCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (ui *UI) ShowError(message string) {
	errorText := ui.styles.error.Render(fmt.Sprintf("❌ 错误: %s", message))
	fmt.Fprintln(os.Stderr, errorText)
}

func (ui *UI) ShowSuccess(message string) {
	successText := ui.styles.success.Render(fmt.Sprintf("✅ %s", message))
	fmt.Println(successText)
}

func (ui *UI) ShowHelp() {
	helpText := `
用法: ccy-core "[失败的命令]" "[错误日志]"
       ccy-core --init
       ccy-core --config

命令:
  --init       输出 Shell 初始化脚本，用于终端集成
  --config     显示配置指南

参数:
  失败的命令  - 执行失败的终端命令
  错误日志    - 可选，终端返回的错误信息

环境变量:
  CCY_API_KEY  - LLM API 密钥 (必需)
  CCY_API_BASE - LLM API 基础 URL (可选, 默认: https://api.openai.com/v1)
  CCY_MODEL    - LLM 模型名称 (可选, 默认: gpt-4)

Shell 集成:
  在 ~/.zshrc 或 ~/.bashrc 中添加: eval "$(ccy-core --init)"
  之后可直接使用 'ccy' 命令，无需手动输入失败的命令和错误

配置:
  开发环境: 在项目根目录创建 .env 文件并设置 CCY_API_KEY
  生产环境: 使用环境变量 export CCY_API_KEY="your-key"
  运行 ccy-core --config 查看详细配置指南
`

	fmt.Println(strings.TrimSpace(helpText))
}

func (ui *UI) ShowConfigGuide() {
	guideText := `
🔧 ccy 配置指南

开发环境 (推荐用于源码开发):
  1. 在项目根目录创建 .env 文件
  2. 添加以下内容:
     CCY_API_KEY=sk-your-api-key-here
     CCY_API_BASE=https://api.openai.com/v1  (可选)
     CCY_MODEL=gpt-4                           (可选)

生产环境 (推荐用于已编译的二进制):
  在 ~/.zshrc 或 ~/.bashrc 中添加:
  export CCY_API_KEY="sk-your-api-key-here"
  export CCY_API_BASE="https://api.openai.com/v1"  (可选)
  export CCY_MODEL="gpt-4"                       (可选)

配置后重启终端或运行:
  source ~/.zshrc  (或 source ~/.bashrc)

验证配置:
  运行任意命令触发 ccy，如配置正确将显示 AI 分析结果
`

	fmt.Println(strings.TrimSpace(guideText))
}
