package common

import (
	"fmt"
	"os"
)

const LOCK_BASE_PATH  = "/tmp/"
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
	fileName := fmt.Sprintf(LOCK_BASE_PATH + "%s.lock", name)
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
	fileName := fmt.Sprintf(LOCK_BASE_PATH + "%s.lock", name)
	if !IsExist(fileName) {
		return false
	}

	if err := os.RemoveAll(fileName); err != nil {
		return false
	}

	return true
}