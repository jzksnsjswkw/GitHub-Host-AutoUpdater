package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

func getHosts() string {
	r, err := http.Get("https://raw.hellogithub.com/hosts")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func main() {
	filePath := `C:\Windows\System32\drivers\etc\hosts`
	h, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	hostFile := string(h)

	reg := regexp.MustCompile(`(?s)# GitHub520 Host Start.*?# GitHub520 Host End\n`)
	b := reg.MatchString(hostFile)
	if b {
		// 存在则更新
		newHost := reg.ReplaceAllString(hostFile, getHosts())
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		f.WriteString(newHost)
	} else {
		// 不存在则追加
		f, err := os.OpenFile(filePath, os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		f.WriteString("\n" + getHosts())
	}
	// 刷新DNS缓存
	err = exec.Command("ipconfig", "/flushdns").Run()
	if err != nil {
		panic(err)
	}
}
