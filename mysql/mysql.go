package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"hello/model"
	"go-src/config"
)

type Mysql struct {
	DB *sql.DB
}

const DRIVE_NAME = "mysql"

var mysqlDB map[string]Mysql
var Conn Mysql
var DSN string

func init() {
	mysqlDB = make(map[string]Mysql)
	Conn = GetInstance()
	DSN = config.DSN
}

func GetInstance() Mysql {
	if mysql, ok := mysqlDB[DSN]; ok {
		return mysql
	}

	db, err := sql.Open(DRIVE_NAME, DSN)
	if err != nil {
		log.Fatal(err)
		return Mysql{}
	}

	mysqlDB[DSN] = Mysql{DB: db}
	return Mysql{DB: db}
}

//单条记录
func (this *Mysql) FindOne(sql string) (model.MapModel, error) {
	rows, err := this.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	//定义输出的类型
	result := make(model.MapModel)
	//这个是sql查询出来的字段
	values := make([]interface{}, count)
	//保存sql查询出来的对应的地址
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		//scansql查询出来的字段的地址
		rows.Scan(valuePtrs...)

		//开始循环columns
		for i, col := range columns {
			var v interface{}
			//值
			val := values[i]
			//判读值的类型（interface类型）如果是byte，则需要转换成字符串
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			//保存
			result[col] = v
		}
	}
	return result, nil
}

//多条记录（根据上面的多条记录修改）
func (this *Mysql) FindAll(sql string) ([]model.MapModel, error) {
	rows, err := this.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]model.MapModel, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(model.MapModel)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	if err != nil {
		return nil, err
	}

	return tableData, nil
}

func (this *Mysql) Exec(sql string) error {
	_, err := this.DB.Exec(sql)
	return err
}