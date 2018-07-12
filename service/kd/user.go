package kd

import (
	"fmt"
	"kd.explorer/tools/mysql"
	"kd.explorer/model"
)

func FindUser(user string) (model.MapModel) {
	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM userinfos WHERE user_key = '%s'", user))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return userInfo
}
