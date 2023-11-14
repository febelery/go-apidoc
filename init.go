package goapidoc

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	tempDir, _ := getTemplateDir()
	destDir := "doc"

	// 判断目录是否存在
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		//目录不存在，执行拷贝操作
		err := copyDir(tempDir, destDir)
		if err != nil {
			slog.Error("init doc error", err)
			return
		}
		slog.Info("init doc successful.")
	}

	// 判断apidoc命令是否存在
	if _, err := exec.LookPath("apidoc"); err != nil {
		cmd := exec.Command("npm", "install", "apidoc", "-g")
		err := cmd.Run()
		if err != nil {
			slog.Error("安装apidoc命令失败：", err)
		} else {
			slog.Info("apidoc 安装成功")
		}
	}
}

// 使用 go list 查询包路径
func getTemplateDir() (string, error) {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/febelery/go-apidoc")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	dir := strings.TrimSpace(string(output)) + "/template"
	return dir, nil
}

// 拷贝文件夹
func copyDir(srcDir, destDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 生成目标路径
		destPath := filepath.Join(destDir, path[len(srcDir):])

		if info.IsDir() {
			// 创建目标文件夹
			err := os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			// 拷贝文件
			err := copyFile(path, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// 拷贝文件
func copyFile(srcFile, destFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}
