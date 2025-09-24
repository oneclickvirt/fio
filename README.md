# fio

一个嵌入fio依赖的golang库(A golang library with embedded fio dependencies)

## 关于 About

这个库提供了对FIO存储性能基准测试工具的Go语言封装。FIO是一个用于测试存储设备性能的综合性基准测试工具，可以测试磁盘的IOPS、带宽和延迟等关键性能指标。

This library provides a Go wrapper for the FIO storage performance benchmark tool. FIO is a comprehensive benchmark tool for testing storage device performance, capable of measuring key performance metrics such as IOPS, bandwidth, and latency of disks.

## 特性 Features

- 支持多平台：Linux (amd64, arm64, 386, arm等), macOS (amd64, arm64), Windows (amd64, 386), FreeBSD (amd64)
- 自动检测并使用系统安装的fio，或使用嵌入的二进制文件
- 自动清理临时文件
- 支持root权限检测和sudo执行
- Multi-platform support: Linux (amd64, arm64, 386, arm, etc.), macOS (amd64, arm64), Windows (amd64, 386), FreeBSD (amd64)
- Automatically detects and uses system-installed fio, or uses embedded binaries
- Automatic cleanup of temporary files
- Support for root permission detection and sudo execution

## 安装 Installation

```bash
go get github.com/oneclickvirt/fio@v0.0.2-20250808045755
```

## 使用方法 Usage

```go
package main

import (
    "log"
    "github.com/oneclickvirt/fio"
)

func main() {
    // 获取fio命令路径
    // Get fio command path
    fioCmd, tempFile, err := fio.GetFIO()
    if err != nil {
        log.Fatalf("Failed to get fio: %v", err)
    }

    // 如果使用了临时文件，确保清理
    // Clean up temporary files if used
    if tempFile != "" {
        defer fio.CleanFio(tempFile)
    }

    // 执行fio基准测试
    // Execute fio benchmark
    err = fio.ExecuteFIO(fioCmd, []string{
        "--name=test",
        "--rw=read",
        "--bs=4k",
        "--size=1G",
        "--numjobs=1",
        "--time_based",
        "--runtime=60s",
        "--group_reporting",
    })
    if err != nil {
        log.Fatalf("Failed to execute fio: %v", err)
    }
}
```

### 简单示例 Simple Example

```go
package main

import (
    "log"
    "github.com/oneclickvirt/fio"
)

func main() {
    // 获取fio命令
    fioCmd, tempFile, err := fio.GetFIO()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // 清理资源
    if tempFile != "" {
        defer fio.CleanFio(tempFile)
    }

    // 执行简单的读取测试
    args := []string{"--name=simple_read_test", "--rw=read", "--bs=4k", "--size=100M"}
    if err := fio.ExecuteFIO(fioCmd, args); err != nil {
        log.Fatalf("FIO execution failed: %v", err)
    }
}
```

## API文档 API Documentation

### 函数 Functions

#### `GetFIO() (string, string, error)`

获取可用的fio命令路径。该函数会首先尝试使用系统安装的fio，如果没有找到则使用嵌入的二进制文件。

Get the available fio command path. This function will first try to use the system-installed fio, and if not found, it will use the embedded binary.

**返回值 Returns:**
- `string`: fio命令的路径或命令字符串 (Path to the fio command or command string)
- `string`: 临时文件路径（如果使用了嵌入的二进制文件）(Temporary file path if using embedded binary)
- `error`: 错误信息 (Error information)

**执行逻辑 Execution Logic:**
1. 检查系统是否已安装fio
2. 如果有root权限，优先尝试`sudo fio`
3. 尝试直接运行fio
4. 如果系统没有fio，则提取嵌入的二进制文件到临时目录

#### `ExecuteFIO(fioCmd string, args []string) error`

执行fio命令。

Execute the fio command.

**参数 Parameters:**
- `fioCmd`: fio命令的路径或字符串 (Path or string for the fio command)
- `args`: 传递给fio的参数 (Arguments to pass to fio)

#### `CleanFio(tempFile string) error`

清理临时文件。

Clean up temporary files.

**参数 Parameters:**
- `tempFile`: 需要清理的临时文件路径 (Path to temporary file to clean up)

## 平台支持 Platform Support

库包含了以下平台的预编译二进制文件：

The library includes precompiled binaries for the following platforms:

### Linux
- amd64 (x86_64)
- 386 (x86 32-bit) 
- arm64 (ARMv8)
- arm (ARMv7)
- riscv64 (RISC-V 64-bit)
- ppc64le (PowerPC64 little-endian)
- ppc64 (PowerPC64 big-endian)
- mips64le (MIPS64 little-endian)
- mips64 (MIPS64 big-endian)
- mipsle (MIPS little-endian)
- mips (MIPS big-endian)
- s390x (IBM System z)

### macOS
- amd64 (Intel)
- arm64 (Apple Silicon)

### Windows
- amd64 (x86_64)
- 386 (x86 32-bit)

### FreeBSD
- amd64 (x86_64)

### 暂缺的二进制文件 Missing Binaries

以下平台的二进制文件暂时不可用，正在开发中：

The following platform binaries are temporarily unavailable and are under development:

- fio-windows-arm64
- fio-freebsd-386  
- fio-freebsd-arm64
- fio-freebsd-arm

## FIO参数示例 FIO Parameter Examples

### 基本测试 Basic Tests

```go
// 顺序读取测试
args := []string{
    "--name=seq_read",
    "--rw=read",
    "--bs=4k", 
    "--size=1G",
    "--numjobs=1",
    "--runtime=60s",
    "--time_based",
}

// 随机写入测试
args := []string{
    "--name=rand_write",
    "--rw=randwrite",
    "--bs=4k",
    "--size=1G", 
    "--numjobs=1",
    "--runtime=60s",
    "--time_based",
}

// 混合读写测试
args := []string{
    "--name=mixed_rw",
    "--rw=randrw",
    "--rwmixread=70",
    "--bs=4k",
    "--size=1G",
    "--numjobs=4",
    "--runtime=60s",
    "--time_based", 
    "--group_reporting",
}
```

### 高级选项 Advanced Options

```go
// IOPS测试
args := []string{
    "--name=iops_test",
    "--rw=randread",
    "--bs=4k",
    "--ioengine=libaio",
    "--iodepth=32", 
    "--direct=1",
    "--size=10G",
    "--numjobs=4",
    "--runtime=300s",
    "--time_based",
    "--group_reporting",
    "--output-format=json",
}
```

## 构建信息 Build Information

该库的二进制文件基于FIO 3.39版本构建。构建脚本和详细信息请参考 `bin/README.md`。

The binaries in this library are built based on FIO version 3.39. For build scripts and detailed information, please refer to `bin/README.md`.

## 许可证 License

GPL-3.0 License