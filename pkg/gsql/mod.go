package gsql

import (
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"database/sql"

	"github.com/chunhui2001/zero4go/pkg/utils"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var Client MySQLClient

type MySQLConf struct {
	Enable   bool   `mapstructure:"MYSQL_ENABLE"`
	Opts     string `mapstructure:"MYSQL_CONN_OPTS" json:"opts"`
	Server   string `mapstructure:"MYSQL_SERVER" json:"server"`
	Database string `mapstructure:"MYSQL_DATABASE" json:"database"`
	User     string `mapstructure:"MYSQL_USER_NAME" json:"user_name"`
	Passwd   string `mapstructure:"MYSQL_PASSWD" json:"passwd"`
	Location string `mapstructure:"MYSQL_MAPPER_LOCATION" json:"mapper_location"`
}

func (c *MySQLConf) connString(passwd string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, c.User, passwd, c.Server, c.Database, c.Opts)
}

var Settings = &MySQLConf{
	Enable:   false,
	Opts:     "timeout=90s&interpolateParams=true&multiStatements=true&charset=utf8&autocommit=true&parseTime=True&loc=Asia%2FShanghai",
	Server:   "127.0.0.1:3306",
	Database: "mydb",
	User:     "keesh",
	Passwd:   "Cc",
}

func Init() {
	if !Settings.Enable {
		Log.Infof("MySQL-Initialized-Disabled: val=%t", Settings.Enable)

		return
	}

	db, err := sql.Open("mysql", Settings.connString(Settings.Passwd))

	if err != nil {
		Log.Errorf("MySQL-Initialized-failed: Error=%s, ConnectionString=%s", err.Error(), Settings.connString("****"))

		return
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		Log.Error(fmt.Sprintf("MySQL-Initialized-failed: Error=%s, ConnectionString=%s", err.Error(), Settings.connString("****")))

		return
	}

	var location = filepath.Join(utils.RootDir(), Settings.Location)

	if tpl, err := template.New("").Funcs(funcMaps()).ParseGlob(location); err != nil {
		panic(err)
	} else {
		Client = MySQLClient{
			DB:     db,
			render: tpl,
		}
	}

	if version, err := Client.Version(); err == nil {
		Log.Info(fmt.Sprintf("MySQL-Initialized-Connected-Successful: ServerVersion=%s, ConnString=%s", version, Settings.connString("****")))
		// execute the Embedding scripts
		//exceScripts()
		return
	}
}
