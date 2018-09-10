package common

import (
	"bytes"
	"errors"
	"fmt"
	"kd.explorer/config"
	"kd.explorer/util/dates"
	"math/rand"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const DefaultSleepTIME = time.Millisecond * 10

func GetLockPath() string {
	path := "/tmp/"
	if IsWindows() {
		path = "C:\\data\\"
	}

	return path
}

func IsWindows() bool {
	return "windows" == runtime.GOOS
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
	return fmt.Sprintf("cd %s;./%s %s", extArr["curDir"], config.TaskList[jobType]["scriptName"], params)
}

/**
 *	获取本机IP
 *	@return string
 */
func GetLocalIp() (string, error) {
	if IsWindows() {
		return "127.0.0.1", nil
	}

	ipData, err := Exec("curl ip.cn")
	if err != nil {
		return "", errors.New("GetLocalIp: " + err.Error())
	}

	reg := regexp.MustCompile(`当前 IP：([0-9.]*).*`)
	localIps := reg.FindStringSubmatch(string(ipData))

	if len(localIps) < 2 {
		return "", errors.New("GetLocalIp: " + strings.Replace(strings.Join(localIps, "|"), "\n", "", -1))
	}

	return localIps[1], nil
}
