package model

import (
	"fmt"
	"kd.explorer/tools/mysql"
)

const UserTable = "userinfos"

func FindUser(user string) (mysql.MapModel) {
	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM %s WHERE user_key = '%s'", UserTable, user))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return userInfo
}
