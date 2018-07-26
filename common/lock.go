package common

import (
	"fmt"
	"os"
	"kd.explorer/config"
	"path"
)

func Lock() bool {
	path := GetLockPath()
	err := MakeDir(path)
	if err != nil {
		panic(err)
	}

	fileName := getLockName()
	if IsExist(fileName) {
		fmt.Println(fileName+" is running")
		return false
	}

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()
	return true
}

func UnLock() bool {
	fileName := getLockName()
	if !IsExist(fileName) {
		fmt.Println(fileName+" is not exists")
		return false
	}

	if err := os.RemoveAll(fileName); err != nil {
		fmt.Println(fileName+" unlock failed")
		return false
	}

	return true
}

func getLockName() string {
	return GetLockPath() + fmt.Sprintf(config.ProNAME + path.Base(os.Args[0])+"_%s.lock", config.CurUser)
}
