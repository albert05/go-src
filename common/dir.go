package common

import (
	"os"
	"path/filepath"
	"strings"
	"kd.explorer/tools/dates"
	"runtime"
	"time"
	"kd.explorer/config"
)

// 检查文件或目录是否存在
// param string filename  文件或目录
// return bool
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 创建目录
// param string dirPath 目录
// return error
func MakeDir(dirPath string) error {
	if !IsExist(dirPath) {
		return os.MkdirAll(dirPath, 0755)
	}

	return nil
}

// 删除文件或空目录
// param string filename 文件或目录
// return error
func RemoveFile(filename string) error {
	if IsExist(filename) {
		return os.Remove(filename)
	}

	return nil
}

// RemoveFilesByPattern 批量删除指定匹配模式的所有文件
func RemoveFilesByPattern(dir string, pattern string) error {
	if IsExist(dir) {
		for _, fileName := range GetAllFileByPattern(dir, pattern) {
			if err := os.RemoveAll(fileName); err != nil {
				return err
			}
		}
	}

	return nil
}

// 删除目录及其内所有文件
// param string dir 文件或目录 （若传文件路径，则删除该文件所在目录）
// return error
func RemoveDir(dir string) error {
	if IsExist(dir) {
		return os.RemoveAll(GetFilePath(dir))
	}

	return nil
}

// 获取指定目录的多级父目录
// param string dir
// param int index 父目录级数
// return string
func GetParentDir(dir string, index int) string {
	if index < 1 {
		return dir
	}
	// 去掉目录末尾的 ‘/’
	dir = strings.TrimRight(dir, "/")

	return GetParentDir(Substr(dir, 0, strings.LastIndex(dir, "/")), index-1)
}

// 获取当前目录
// return string
func GetPwd() string {
	dir, err := os.Getwd()
	if err == nil {
		return strings.Replace(dir, "\\", "/", -1)
	}

	return ""
}

// 获取程序根目录
// return string
func GetRootDir() string {
	return GetParentDir(GetPwd(), 1)
}

// 获取文件目录
// param string file
// return string
func GetFilePath(file string) string {
	if stat, err := os.Stat(file); err == nil && stat.IsDir() {
		return file
	}

	return filepath.Dir(file)
}

// 获取目录所有文件
// param string dir
// return slice
func GetAllFilesByDir(dir string) []string {
	pattern := GetFilePath(dir) + "/*"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	return files
}

// 获取目录指定后缀名的所有文件
// param string dir
// param string suffix 文件后缀 如 .txt
// return slice
func GetAllFileByDirSuffix(dir string, suffix string) []string {
	pattern := GetFilePath(dir) + GetSysDirSuffix() + "*" + suffix
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	return files
}

func GetSysDirSuffix() string {
	suffix := "/"
	if "windows" == runtime.GOOS {
		suffix = "\\"
	}

	return suffix
}

// 获取目录指定匹配模式的所有文件
// param string dir
// param string pattern 匹配模式
// return slice
func GetAllFileByPattern(dir string, pattern string) []string {
	filePattern := GetFilePath(dir) + GetSysDirSuffix() + pattern
	files, err := filepath.Glob(filePattern)
	if err != nil {
		return nil
	}

	return files
}

func GetLogPath(jobType string) string {
	path := config.LogPATH + jobType + "/"
	if err := MakeDir(path); err != nil {
		return ""
	}

	return " 1> " + path + dates.NowDateShortStr() + ".log 2>&1"
}

func GetFileModTime(file string) int64 {
	f, err := os.Open(file)
	if err != nil {
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}