package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ccy-ai/ccy-assistant/internal/api"
	"github.com/ccy-ai/ccy-assistant/internal/config"
	"github.com/ccy-ai/ccy-assistant/internal/model"
	"github.com/ccy-ai/ccy-assistant/internal/shellinit"
	"github.com/ccy-ai/ccy-assistant/internal/tui"
)

func main() {
	args := os.Args[1:]

	ui := tui.NewUI()

	if len(args) == 0 {
		ui.ShowHelp()
		os.Exit(0)
	}

	if args[0] == "--init" {
		script, err := shellinit.GenerateInitScript()
		if err != nil {
			ui.ShowError(err.Error())
			os.Exit(1)
		}
		fmt.Println(script)
		os.Exit(0)
	}

	if args[0] == "--config" {
		ui.ShowConfigGuide()
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		ui.ShowError("🔴 未检测到 API Key。请在 .env 文件中配置，或在终端执行 export CCY_API_KEY='你的密钥'。")
		ui.ShowConfigGuide()
		os.Exit(1)
	}

	failedCommand := args[0]
	var errorMessage string
	if len(args) > 1 {
		errorMessage = args[1]
	}

	client := api.NewClient(cfg)

	stopLoading := ui.ShowLoading("正在呼叫 AI 专家")
	responseText, err := client.SendRequest(failedCommand, errorMessage)
	stopLoading()

	if err != nil {
		ui.ShowError(err.Error())
		os.Exit(1)
	}

	var response model.AnalysisResponse
	if err := json.Unmarshal([]byte(responseText), &response); err != nil {
		ui.ShowError(fmt.Sprintf("解析 API 响应失败: %v", err))
		fmt.Printf("\n原始响应:\n%s\n", responseText)
		os.Exit(1)
	}

	if response.Command == "" {
		ui.ShowError("API 返回的命令为空")
		os.Exit(1)
	}

	ui.ShowResult(response.Analysis, response.Command)

	confirmed, err := ui.AskConfirmation()
	if err != nil {
		ui.ShowError(fmt.Sprintf("读取输入失败: %v", err))
		os.Exit(1)
	}

	if confirmed {
		if err := ui.ExecuteCommand(response.Command); err != nil {
			ui.ShowError(fmt.Sprintf("执行命令失败: %v", err))
			os.Exit(1)
		}
		ui.ShowSuccess("命令执行完成")
	} else {
		fmt.Println("\n已取消")
	}
}
