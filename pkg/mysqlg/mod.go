package mysqlg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thoas/go-funk"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var MysqlClient *sql.DB

type MySqlConf struct {
	Enable   bool   `mapstructure:"MYSQL_ENABLE"`
	Opts     string `mapstructure:"MYSQL_CONN_OPTS" json:"opts"`
	Server   string `mapstructure:"MYSQL_SERVER" json:"server"`
	Database string `mapstructure:"MYSQL_DATABASE" json:"database"`
	User     string `mapstructure:"MYSQL_USER_NAME" json:"user_name"`
	Passwd   string `mapstructure:"MYSQL_PASSWD" json:"passwd"`
}

func (c *MySqlConf) connString(passwd string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, c.User, passwd, c.Server, c.Database, c.Opts)
	//return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, c.User, passwd, c.Server, c.Database)
}

var MysqlSetting = &MySqlConf{
	Enable:   false,
	Opts:     "timeout=90s&interpolateParams=true&multiStatements=true&charset=utf8&autocommit=true&parseTime=True&loc=Asia%2FShanghai",
	Server:   "127.0.0.1:3306",
	Database: "mydb",
	User:     "keesh",
	Passwd:   "Cc",
}

func Init() {
	if !MysqlSetting.Enable {
		Log.Infof("Init mysql mod: val=%b", MysqlSetting.Enable)
		return
	}

	db, err := sql.Open("mysql", MysqlSetting.connString(MysqlSetting.Passwd))

	if err != nil {
		Log.Errorf("mysql init failed: Error=%s, ConnectionString=%s", err.Error(), MysqlSetting.connString("****"))

		return
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		Log.Error(fmt.Sprintf("Mysql-Client-Error: Error=%s, ConnectionString=%s", err.Error(), MysqlSetting.connString("****")))

		return
	}

	MysqlClient = db

	if version, err := Version(); err == nil {
		Log.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: ServerVersion=%s, ConnString=%s", version, MysqlSetting.connString("****")))
		// execute the Embedding scripts
		//exceScripts()
		return
	}
}

func Version() (string, error) {
	var version string
	err2 := MysqlClient.QueryRow("SELECT VERSION()").Scan(&version)

	return version, err2
}

func Exec(sqlStr string, timeout int32, args ...any) (sql.Result, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond*1000000000*time.Duration(timeout))

	defer cancelFunc()

	// prepare the statement
	stmt, err := MysqlClient.Prepare(sqlStr)

	if err != nil {
		theSql := sqlStr

		Log.Error(fmt.Sprintf("Mysql-Insert-Error-1: sqlStr=%s, Error=%s", theSql, err.Error()))

		return nil, err
	}

	// format all args at once
	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		Log.Error(fmt.Sprintf("Mysql-Insert-Error-2: sqlStr=%s, Error=%s", sqlStr, err.Error()))

		return nil, err
	}

	return result, nil
}

// 批量插入
// INSERT INTO tableName(id, name) VALUES (?, ?),(?, ?),(?, ?)
func InsertBulk(timeout int32, tableName string, columeMaps [][]string, insertData []map[string]interface{}) (sql.Result, error) {
	_columes := funk.Map(columeMaps, func(m []string) string { return m[0] }).([]string)
	_keys := funk.Map(columeMaps, func(m []string) string { return m[1] }).([]string)

	insert := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(_columes[:], ", "))

	vals := make([]any, 0, len(insertData))
	placeholders := strings.Repeat("?, ", len(_columes))

	for _, item := range insertData {
		insert += fmt.Sprintf(`(%s),`, placeholders[:len(placeholders)-2])

		for _, _k := range _keys {
			if item[_k] == nil {
				vals = append(vals, nil)
			} else {
				vals = append(vals, item[_k].(interface{}))
			}
		}
	}

	return Exec(insert[:len(insert)-1]+";", timeout, vals...)
}
