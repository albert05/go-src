package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"kd.explorer/config"
	"strings"
	"unsafe"
)

type Mysql struct {
	DB *sql.DB
}

type queryResult struct {
	affectedRows, insertId int64
}

const DriverNAME = "mysql"

var mysqlDB map[string]Mysql
var Conn Mysql
var DSN string

func init() {
	DSN = config.DSN
	mysqlDB = make(map[string]Mysql)
	Conn = GetInstance()
}

func GetInstance() Mysql {
	if mysql, ok := mysqlDB[DSN]; ok {
		return mysql
	}

	db, err := sql.Open(DriverNAME, DSN)
	if err != nil {
		panic(err)
		return Mysql{}
	}

	mysqlDB[DSN] = Mysql{DB: db}
	return Mysql{DB: db}
}

//单条记录
func (this *Mysql) FindOne(sql string) (MapModel, error) {
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
	result := make(MapModel)
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
func (this *Mysql) FindAll(sql string) ([]MapModel, error) {
	rows, err := this.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]MapModel, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(MapModel)
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

// 插入数据
// string table
// map data 插入的数据
// return bool
func (this *Mysql) Insert(table string, data map[string]string) bool {
	if data == nil {
		return false
	}
	var length int = len(data)
	var i int
	var columns = make([]string, length)
	values := make([]interface{}, length)

	for key, value := range data {
		columns[i] = key
		values[i] = value
		i++
	}
	columnStr := strings.Join(columns, "`, `")
	columnRep := strings.Repeat("?,", length-1)
	sql := "INSERT INTO " + table + " (`" + columnStr + "`) VALUES (" + columnRep + "?)"

	stmtIns, err := this.DB.Prepare(sql) // ? = placeholder
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer stmtIns.Close()

	result, err := stmtIns.Exec(values...) // Insert tuples (i, i^2)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var qResult queryResult = *(*queryResult)(unsafe.Pointer(&result))
	return qResult.affectedRows > 0
}

// 更新数据
// string table
// map data 插入的数据
// map condition 更新条件
// return bool
func (this *Mysql) Update(table string, data map[string]string, condition map[string]string) bool {
	var i int
	length := len(data)
	columns := make([]string, length)
	values := make([]interface{}, length)

	for key, value := range data {
		columns[i] = key + "=?"
		values[i] = value
		i++
	}
	where, _ := condition["where"]

	sql := "UPDATE " + table + " SET " + strings.Join(columns, ",") + " WHERE " + where

	stmtIns, err := this.DB.Prepare(sql) // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()

	result, err := stmtIns.Exec(values...) // Insert tuples (i, i^2)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var qResult queryResult = *(*queryResult)(unsafe.Pointer(&result))
	return qResult.affectedRows > 0
}

func (this *Mysql) Exec(sql string) error {
	_, err := this.DB.Exec(sql)
	return err
}
