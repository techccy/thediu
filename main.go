package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ccy-ai/ccy-assistant/internal/config"
	"github.com/ccy-ai/ccy-assistant/internal/context"
	"github.com/ccy-ai/ccy-assistant/internal/memory"
	"github.com/ccy-ai/ccy-assistant/internal/provider"
	"github.com/ccy-ai/ccy-assistant/internal/shellinit"
	"github.com/ccy-ai/ccy-assistant/internal/tui"
)

type AnalysisResponse struct {
	Analysis string `json:"analysis"`
	Command  string `json:"command"`
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		ui := tui.NewUI()
		ui.ShowHelp()
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		ui := tui.NewUI()
		ui.ShowError(fmt.Sprintf("配置加载失败: %v", err))
		os.Exit(1)
	}

	ui := tui.NewUIWithConfig(cfg)

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

	if args[0] == "--switch" {
		if err := ui.ShowProviderSwitch(); err != nil {
			ui.ShowError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if args[0] == "--help" || args[0] == "-h" {
		ui.ShowHelp()
		os.Exit(0)
	}

	failedCommand := args[0]
	var errorMessage string
	if len(args) > 1 {
		errorMessage = args[1]
	}

	mem, err := memory.NewMemory()
	if err != nil {
		ui.ShowError(fmt.Sprintf("记忆模块初始化失败: %v", err))
		os.Exit(1)
	}
	defer mem.Close()

	cachedEntry, err := mem.Find(failedCommand, errorMessage)
	if err == nil && cachedEntry != nil {
		ui.ShowResult(fmt.Sprintf("[⚡️ Local Cache] %s", cachedEntry.Error), cachedEntry.FixCommand)
		if confirmed, _ := ui.AskConfirmation(); confirmed {
			if err := ui.ExecuteCommand(cachedEntry.FixCommand); err != nil {
				ui.ShowError(fmt.Sprintf("执行命令失败: %v", err))
			} else {
				ui.ShowSuccess("命令执行完成")
			}
		}
		os.Exit(0)
	}

	providerName, _, err := cfg.GetCurrentProvider()
	if err != nil {
		ui.ShowError(fmt.Sprintf("获取提供商配置失败: %v", err))
		os.Exit(1)
	}

	factory := provider.NewProviderFactory()
	llmProvider, err := factory.CreateProvider(cfg, providerName)
	if err != nil {
		ui.ShowError(fmt.Sprintf("创建提供商失败: %v", err))
		os.Exit(1)
	}

	ctxEngine := context.NewContext()
	workDir, _ := os.Getwd()
	additionalContext, _ := ctxEngine.Collect(workDir, errorMessage)

	systemPrompt := `You are a Linux/Unix system expert. Your task is to analyze the failed command and error message, then provide a corrected command.

You must respond with a valid JSON object only, without any markdown code blocks or extra text. The JSON must have this exact format:
{
  "analysis": "Brief explanation of the error cause in one sentence",
  "command": "The corrected command that should be executed"
}

Example:
{
  "analysis": "The command 'pushu' is not a valid git subcommand",
  "command": "git push origin main"
}`

	var userContent string
	if errorMessage != "" {
		userContent = fmt.Sprintf("Failed command: %s\nError message: %s", failedCommand, errorMessage)
	} else {
		userContent = fmt.Sprintf("Failed command: %s", failedCommand)
	}

	if additionalContext != "" {
		userContent += additionalContext
	}

	messages := []provider.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	stopLoading := ui.ShowLoading(fmt.Sprintf("正在呼叫 %s (%s)", providerName, llmProvider.GetModelName()))
	var responseText string
	var requestErr error

	responseText, requestErr = llmProvider.SendRequest(messages)

	if requestErr != nil && providerName != "ollama" {
		stopLoading()
		ui.ShowSuccess(fmt.Sprintf("网络请求失败，尝试切换到本地 Ollama..."))

		if ollamaProvider, err := factory.CreateProvider(cfg, "ollama"); err == nil {
			stopLoading = ui.ShowLoading(fmt.Sprintf("正在呼叫 Ollama (%s)", ollamaProvider.GetModelName()))
			responseText, requestErr = ollamaProvider.SendRequest(messages)
		}
	}

	stopLoading()

	if requestErr != nil {
		ui.ShowError(fmt.Sprintf("请求失败: %v", requestErr))
		os.Exit(1)
	}

	var response AnalysisResponse
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
			mem.Save(failedCommand, errorMessage, response.Command, false)
			ui.ShowError(fmt.Sprintf("执行命令失败: %v", err))
			os.Exit(1)
		}
		mem.Save(failedCommand, errorMessage, response.Command, true)
		ui.ShowSuccess("命令执行完成")
	} else {
		fmt.Println("\n已取消")
	}
}
