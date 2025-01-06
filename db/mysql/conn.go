package mysql

import (
	"cloud_storage/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", config.MySQLSource)
	fmt.Printf("初始化db\n")
	if err != nil {
		fmt.Printf("连接失败", err.Error())
		os.Exit(1)
	}
	db.SetMaxOpenConns(1000)
	db.SetConnMaxLifetime(time.Second * 300)
	err = db.Ping()
	if err != nil {
		fmt.Printf("连接失败2", err.Error())
		os.Exit(1)
	}
}

// DBConn : 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
