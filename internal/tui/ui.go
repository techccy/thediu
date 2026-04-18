package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type UI struct {
	reader *bufio.Reader
}

func NewUI() *UI {
	return &UI{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (ui *UI) ShowLoading(message string) func() {
	fmt.Printf("\r%s...", message)
	return func() {
		fmt.Print("\r")
	}
}

func (ui *UI) ShowResult(analysis, command string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("\x1b[36m分析:\x1b[0m %s\n", analysis)
	fmt.Printf("\x1b[33m建议命令:\x1b[0m \x1b[32m%s\x1b[0m\n", command)
	fmt.Println(strings.Repeat("=", 60))
}

func (ui *UI) AskConfirmation() (bool, error) {
	fmt.Print("\n\x1b[36m执行此命令? \x1b[33m[Y/n]\x1b[0m: ")

	input, err := ui.reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes", nil
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
	fmt.Fprintf(os.Stderr, "\x1b[31m错误: %s\x1b[0m\n", message)
}

func (ui *UI) ShowHelp() {
	fmt.Println("用法: ccy-core \"[失败的命令]\" \"[错误日志]\"")
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  失败的命令  - 执行失败的终端命令")
	fmt.Println("  错误日志    - 可选，终端返回的错误信息")
	fmt.Println()
	fmt.Println("环境变量:")
	fmt.Println("  CCY_API_KEY  - LLM API 密钥 (必需)")
	fmt.Println("  CCY_API_BASE - LLM API 基础 URL (可选, 默认: https://api.openai.com/v1)")
	fmt.Println("  CCY_MODEL    - LLM 模型名称 (可选, 默认: gpt-4)")
}
