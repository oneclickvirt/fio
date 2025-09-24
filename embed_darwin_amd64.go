//go:build darwin && amd64
// +build darwin,amd64

package fio

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed bin/fio-darwin-amd64
var binFiles embed.FS

// GetFIO 返回适用于当前系统的 fio 命令字符串（带不带 sudo）和临时文件路径（用于清理）
func GetFIO() (fioCmd string, tempFile string, err error) {
	var errors []string
	// 1. 尝试系统自带 fio
	if path, lookErr := exec.LookPath("fio"); lookErr == nil {
		if hasRootPermission() {
			// 尝试 sudo fio
			testCmd := exec.Command("sudo", path, "--help")
			if runErr := testCmd.Run(); runErr == nil {
				return "sudo fio", "", nil
			} else {
				errors = append(errors, fmt.Sprintf("sudo fio 测试失败: %v", runErr))
			}
		}
		// 直接尝试 fio
		testCmd := exec.Command(path, "--help")
		if runErr := testCmd.Run(); runErr == nil {
			return "fio", "", nil
		} else {
			errors = append(errors, fmt.Sprintf("fio 直接运行失败: %v", runErr))
		}
	} else {
		errors = append(errors, fmt.Sprintf("无法找到 fio: %v", lookErr))
	}
	// 2. 创建临时目录
	tempDir, tempErr := os.MkdirTemp("", "fioWrapper")
	if tempErr != nil {
		return "", "", fmt.Errorf("创建临时目录失败: %v", tempErr)
	}
	// 3. 尝试使用嵌入的 fio 版本
	binName := "fio-darwin-amd64"
	binPath := filepath.Join("bin", binName)
	fileContent, readErr := binFiles.ReadFile(binPath)
	if readErr == nil {
		tempFile = filepath.Join(tempDir, binName)
		writeErr := os.WriteFile(tempFile, fileContent, 0755)
		if writeErr == nil {
			if hasRootPermission() {
				// 尝试 sudo 嵌入版本
				testCmd := exec.Command("sudo", tempFile, "--help")
				if runErr := testCmd.Run(); runErr == nil {
					return fmt.Sprintf("sudo %s", tempFile), tempFile, nil
				} else {
					errors = append(errors, fmt.Sprintf("sudo %s 运行失败: %v", tempFile, runErr))
				}
			}
			// 直接尝试嵌入版本
			testCmd := exec.Command(tempFile, "--help")
			if runErr := testCmd.Run(); runErr == nil {
				return tempFile, tempFile, nil
			} else {
				errors = append(errors, fmt.Sprintf("%s 运行失败: %v", tempFile, runErr))
			}
		} else {
			errors = append(errors, fmt.Sprintf("写入临时文件失败 (%s): %v", tempFile, writeErr))
		}
	} else {
		errors = append(errors, fmt.Sprintf("读取嵌入的 fio 二进制文件失败: %v", readErr))
	}
	// 返回所有错误信息
	return "", "", fmt.Errorf("无法找到可用的 fio 命令:\n%s", strings.Join(errors, "\n"))
}

// ExecuteFIO 执行拼好的 fio 命令字符串（包括 sudo、fio 等）
func ExecuteFIO(fioCmd string, args []string) error {
	// 拼接命令字符串
	fullCmd := fmt.Sprintf("%s %s", fioCmd, strings.Join(args, " "))
	cmd := exec.Command("sh", "-c", fullCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
