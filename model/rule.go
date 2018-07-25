package model

import (
	"fmt"
	"kd.explorer/util/mysql"
)

const RuleTable = "rule_list"

func FindRule(key string) mysql.MapModel {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE k='%s'", RuleTable, key)
	rule, err := mysql.Conn.FindOne(sql)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return rule
}
