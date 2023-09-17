package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// 实时输出命令的结果
func HandleCommandOutput(cmd *exec.Cmd) (string, error) {
	var output strings.Builder

	// 创建管道
	pr, pw := io.Pipe()
	cmd.Stdout = pw // 将命令的标准输出重定向到管道的写入端
	cmd.Stderr = pw // 将命令的标准错误输出重定向到管道的写入端
	fmt.Printf("开始执行命令:%v\n", cmd.String())
	// 实时输出命令的结果
	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line)
			output.WriteString("\n")
			fmt.Println(line)
		}
	}()

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	pw.Close()

	return output.String(), nil
}
