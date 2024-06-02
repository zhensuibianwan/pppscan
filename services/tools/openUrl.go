package tools

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"syscall"
)

func OpenUrlByDefaultBrowser(rawURL string) {
	// 验证URL的格式
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("无效的URL")
		return
	}

	// 验证URL的Scheme是否为http或https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		fmt.Println("URL必须以http或https开头")
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "/b", rawURL)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}
