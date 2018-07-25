package model

import (
	"fmt"
	"kd.explorer/util/dates"
	"kd.explorer/util/mysql"
)

const MonitorTable = "trans_monitor_list"

func FindMRecord(transId string) mysql.MapModel {
	monitorRecord, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM %s WHERE trans_id = '%s'", MonitorTable, transId))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return monitorRecord
}

func InsertMRecord(transId, data string) bool {
	datas := map[string]string{
		"trans_id":   transId,
		"data":       data,
		"created_at": fmt.Sprintf("%d", dates.NowTime()),
	}

	return mysql.Conn.Insert(MonitorTable, datas)
}
