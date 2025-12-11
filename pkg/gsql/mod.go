package gsql

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/flosch/pongo2/v6"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var MysqlClient MySqlClient

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

	MysqlClient = MySqlClient{
		DB:     db,
		render: pongo2.NewSet("my_templates", pongo2.MustNewLocalFileSystemLoader("./META-INF/mappers")),
	}

	// 注册过滤器
	RegisterSqlBindFilter()

	if version, err := MysqlClient.Version(); err == nil {
		Log.Info(fmt.Sprintf("Mysql-Client-Connected-Successful: ServerVersion=%s, ConnString=%s", version, MysqlSetting.connString("****")))
		// execute the Embedding scripts
		//exceScripts()
		return
	}
}
