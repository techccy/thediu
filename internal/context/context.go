package context

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxContextSize = 10 * 1024
	maxDepth       = 2
)

var sensitiveFiles = []string{
	".pem", ".key", ".env", ".env.local", ".env.production",
	"id_rsa", "id_ed25519", "id_dsa", ".ssh",
	"credentials.json", "service-account.json",
}

type Context struct {
	DirectoryStructure string
	FileContents       map[string]string
}

func NewContext() *Context {
	return &Context{
		FileContents: make(map[string]string),
	}
}

func (c *Context) Collect(workDir string, errorMessage string) (string, error) {
	var builder strings.Builder

	if dirStructure, err := c.getDirectoryStructure(workDir); err == nil {
		c.DirectoryStructure = dirStructure
		builder.WriteString("\n=== Directory Structure ===\n")
		builder.WriteString(dirStructure)
		builder.WriteString("\n")
	}

	if err := c.sniffFileContents(errorMessage, workDir); err == nil {
		if len(c.FileContents) > 0 {
			builder.WriteString("\n=== Relevant File Contents ===\n")
			for filePath, content := range c.FileContents {
				builder.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", filePath, content))
			}
		}
	}

	return builder.String(), nil
}

func (c *Context) getDirectoryStructure(dir string) (string, error) {
	var result strings.Builder
	ignoreDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true,
		"target": true, "build": true, "dist": true,
		".idea": true, ".vscode": true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return nil
		}

		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) > maxDepth {
			return filepath.SkipDir
		}

		if info.IsDir() {
			if ignoreDirs[info.Name()] {
				return filepath.SkipDir
			}
			result.WriteString(fmt.Sprintf("%s/\n", relPath))
		} else {
			result.WriteString(fmt.Sprintf("%s\n", relPath))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (c *Context) sniffFileContents(errorMessage, workDir string) error {
	patterns := []string{
		`([a-zA-Z0-9_\-./]+\.go):(\d+)`,
		`([a-zA-Z0-9_\-./]+\.py):(\d+)`,
		`([a-zA-Z0-9_\-./]+\.js):(\d+)`,
		`([a-zA-Z0-9_\-./]+\.ts):(\d+)`,
		`([a-zA-Z0-9_\-./]+\.java):(\d+)`,
		`([a-zA-Z0-9_\-./]+\.rs):(\d+)`,
		`File "([a-zA-Z0-9_\-./]+)", line (\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(errorMessage, -1)

		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			filePath := match[1]
			lineNum := match[2]

			if c.isSensitive(filePath) {
				continue
			}

			if c.getTotalSize() >= maxContextSize {
				continue
			}

			fullPath := filepath.Join(workDir, filePath)
			if content, err := c.readFileContext(fullPath, lineNum); err == nil {
				c.FileContents[filePath] = content
			}
		}
	}

	return nil
}

func (c *Context) readFileContext(filePath, lineNumStr string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	lineNum := 1
	fmt.Sscanf(lineNumStr, "%d", &lineNum)

	start := lineNum - 11
	if start < 0 {
		start = 0
	}

	end := lineNum + 10
	if end > len(lines) {
		end = len(lines)
	}

	var result strings.Builder
	for i := start; i < end; i++ {
		prefix := "  "
		if i == lineNum-1 {
			prefix = "> "
		}
		result.WriteString(fmt.Sprintf("%s%4d: %s\n", prefix, i+1, lines[i]))
	}

	return result.String(), nil
}

func (c *Context) isSensitive(filePath string) bool {
	base := filepath.Base(filePath)
	for _, suffix := range sensitiveFiles {
		if strings.HasSuffix(base, suffix) || strings.Contains(base, suffix) {
			return true
		}
	}
	return false
}

func (c *Context) getTotalSize() int {
	total := 0
	for _, content := range c.FileContents {
		total += len(content)
	}
	return total
}
