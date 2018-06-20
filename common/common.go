package common

import (
	"fmt"
	"os"
	"math/rand"
	"os/exec"
	"bytes"
)

const LockBasePATH  = "/tmp/"
var Config map[string]string

func initConfig(path string) {

}

func GetConfig(key string) {

}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if pos > len(runes) {
		return ""
	}
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Lock(name string) bool {
	fileName := fmt.Sprintf(LockBasePATH + "%s.lock", name)
	if IsExist(fileName) {
		return false
	}

	f, err := os.Create(fileName)
	if  err != nil {
		return false
	}

	defer f.Close()
	return true
}

func UnLock(name string) bool {
	fileName := fmt.Sprintf(LockBasePATH + "%s.lock", name)
	if !IsExist(fileName) {
		return false
	}

	if err := os.RemoveAll(fileName); err != nil {
		return false
	}

	return true
}

func GenerateRangeNum(min, max int) int {
	randNum := rand.Intn(max - min) + min
	return randNum
}


// 同步执行命令
func Cmd(cmdStr string) error {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Start()
	status := cmd.Wait()
	return status
}

// 同步执行命令, 并返回执行的结果
func Exec(cmdStr string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}
