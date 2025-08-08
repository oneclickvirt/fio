//go:build windows && 386
// +build windows,386

package fio

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed bin/fio-windows-386.exe
var binFiles embed.FS

// GetFIO 返回适用于当前系统的 fio 命令字符串和临时文件路径（用于清理）
func GetFIO() (fioCmd string, tempFile string, err error) {
	var errors []string
	// 1. 尝试系统自带 fio
	if path, lookErr := exec.LookPath("fio.exe"); lookErr == nil {
		// 直接尝试 fio.exe (Windows 通常不使用 sudo)
		testCmd := exec.Command(path, "--help")
		if runErr := testCmd.Run(); runErr == nil {
			return path, "", nil
		} else {
			errors = append(errors, fmt.Sprintf("fio.exe 运行失败: %v", runErr))
		}
	} else if path, lookErr := exec.LookPath("fio"); lookErr == nil {
		// 尝试不带 .exe 后缀
		testCmd := exec.Command(path, "--help")
		if runErr := testCmd.Run(); runErr == nil {
			return path, "", nil
		} else {
			errors = append(errors, fmt.Sprintf("fio 运行失败: %v", runErr))
		}
	} else {
		errors = append(errors, fmt.Sprintf("无法找到 fio: %v", lookErr))
	}
	// 2. 创建临时目录
	tempDir, tempErr := os.MkdirTemp("", "fioWrapper")
	if tempErr != nil {
		return "", "", fmt.Errorf("创建临时目录失败: %v", tempErr)
	}
	// 3. 使用嵌入的 fio 版本
	binName := "fio-windows-386.exe"
	binPath := fmt.Sprintf("bin/%s", binName)
	fileContent, readErr := binFiles.ReadFile(binPath)
	if readErr == nil {
		tempFile = filepath.Join(tempDir, binName)
		writeErr := os.WriteFile(tempFile, fileContent, 0755)
		if writeErr == nil {
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

// ExecuteFIO 执行拼好的 fio 命令字符串
func ExecuteFIO(fioCmd string, args []string) error {
	// Windows 版本不使用 sh -c 方式执行，而是直接执行命令
	cmd := exec.Command(fioCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
