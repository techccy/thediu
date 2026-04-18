package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/techccy/diu-assistant/internal/config"
	"github.com/techccy/diu-assistant/internal/context"
	"github.com/techccy/diu-assistant/internal/provider"
	"github.com/techccy/diu-assistant/internal/shellinit"
	"github.com/techccy/diu-assistant/internal/tui"
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

	systemPrompt := `You are an elite developer and system administrator. Your task is to analyze a failed terminal command and its error output, then provide the most direct command to fix the issue.

CRITICAL RULES:
1. Shell Typos: If it's a misspelled command, correct it.
2. System Tools: If a CLI tool is missing, provide the install command (e.g., apt/brew).
3. Missing Dependencies (Programming Languages): If the error is a missing library (e.g., Python 'ModuleNotFoundError', Node 'Cannot find module'), your command MUST be the package manager installation command (e.g., 'pip install <pkg>', 'npm install <pkg>'). 
   * IMPORTANT: Be smart about package names. For example, in Python, 'cv2' is 'opencv-python', 'PIL' is 'Pillow', 'yaml' is 'PyYAML'.
4. Chaining: You can use '&&' to chain the fix and re-run the original command if appropriate.

You must respond with a valid JSON object ONLY, without any markdown code blocks (no ` + "```json" + `) or extra text. 

FORMAT:
{
  "analysis": "Brief explanation of the error cause in one sentence (in Chinese)",
  "command": "The exact terminal command to execute"
}

EXAMPLES:
Input command: "git pushu"
Input error: "git: 'pushu' is not a git command"
Output: {"analysis": "Git命令拼写错误，应为push", "command": "git push"}

Input command: "python main.py"
Input error: "ModuleNotFoundError: No module named 'requests'"
Output: {"analysis": "Python运行缺少requests依赖包", "command": "pip install requests && python main.py"}

Input command: "python script.py"
Input error: "ModuleNotFoundError: No module named 'cv2'"
Output: {"analysis": "缺少OpenCV库，Python中对应的包名为opencv-python", "command": "pip install opencv-python && python script.py"}`

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
		ui.ShowSuccess("网络请求失败，尝试切换到本地 Ollama...")

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
			ui.ShowError(fmt.Sprintf("执行命令失败: %v", err))
			os.Exit(1)
		}
		ui.ShowSuccess("命令执行完成")
	} else {
		fmt.Println("\n已取消")
	}
}
