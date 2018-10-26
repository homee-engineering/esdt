package io

import (
	"fmt"
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

var box = packr.NewBox("../../static")

func GetFile(filePath string) (string, error) {
	b := box.Bytes(filePath)
	file := filepath.Base(filePath)
	err := ioutil.WriteFile(file, b, os.ModePerm)
	if err != nil {
		return "", err
	}
	return file, nil
}

func ExecCommand(c *exec.Cmd) ([]byte, error) {
	stdin, err := c.CombinedOutput()

	if err != nil {
		return stdin, err
	}
	return stdin, nil
}

func RunCommand(s string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", s)
	return ExecCommand(cmd)
}

func Cmd(s string, d ...interface{}) ([]byte, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(s, d...))
	return ExecCommand(cmd)
}

func ApplyTemplate(fileName string, data interface{}) (string, error) {
	s := box.String(fileName)

	t, err := template.New("").Parse(s)

	if err != nil {
		return "", err
	}

	filePath := filepath.Join(os.TempDir(), fileName)
	fo, err := os.Create(filePath)
	err = t.Execute(fo, data)

	if err != nil {
		return "", err
	}

	return filePath, nil
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func Cwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
