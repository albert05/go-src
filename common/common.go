package common

import (
	"bytes"
	"errors"
	"fmt"
	"kd.explorer/config"
	"kd.explorer/util/dates"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const DefaultSleepTIME = time.Millisecond * 10

func GetLockPath() string {
	path := "/tmp/"
	if "windows" == runtime.GOOS {
		path = "C:\\data\\"
	}

	return path
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
	path := GetLockPath()
	err := MakeDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf(path+"%s.lock", name)
	if IsExist(fileName) {
		return false
	}

	f, err := os.Create(fileName)
	if err != nil {
		return false
	}

	defer f.Close()
	return true
}

func UnLock(name string) bool {
	fileName := fmt.Sprintf(GetLockPath()+"%s.lock", name)
	if !IsExist(fileName) {
		return false
	}

	if err := os.RemoveAll(fileName); err != nil {
		return false
	}

	return true
}

func GenerateRangeNum(min, max int) int {
	randNum := rand.Intn(max-min) + min
	return randNum
}

// 异步执行命令
func Cmd(cmdStr string) {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Start()
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

func Wait(timePoint float64) {
	currTime := dates.TimeInt2float(dates.CurrentMicro())
	fmt.Println(currTime, timePoint)

	for currTime < timePoint {
		time.Sleep(DefaultSleepTIME)

		currTime = dates.TimeInt2float(dates.CurrentMicro())
	}
}

func GetCmdStr(jobType string, extArr map[string]string) string {
	params := fmt.Sprintf(config.TaskList[jobType]["params"], jobType, extArr["ids"])
	return fmt.Sprintf("cd %s;./%s %s %s", extArr["curDir"], config.TaskList[jobType]["scriptName"], params, extArr["logDir"])
}

/**
 *	获取本机IP
 *	@return string
 */
func GetLocalIp() (string, error) {
	ipData, err := Exec("curl ip.cn")
	if err != nil {
		return "", errors.New("GetLocalIp: " + err.Error())
	}

	reg := regexp.MustCompile(`当前 IP：([0-9.]*).*`)
	localIps := reg.FindStringSubmatch(string(ipData))
	fmt.Println(localIps)

	if len(localIps) < 2 {
		return "", errors.New("GetLocalIp: " + strings.Replace(strings.Join(localIps, "|"), "\n", "", -1))
	}

	return localIps[1], nil
}
