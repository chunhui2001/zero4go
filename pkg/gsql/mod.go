package gsql

import (
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"database/sql"

	"github.com/chunhui2001/zero4go/pkg/utils"
	_ "github.com/go-sql-driver/mysql"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

var Client MySQLClient
var DataSouces map[string]*MySQLClient

type MySQLConf struct {
	Enable   bool   `mapstructure:"MYSQL_ENABLE"`
	Name     string `mapstructure:"MYSQL_NAME"`
	Server   string `mapstructure:"MYSQL_SERVER" json:"server"`
	Database string `mapstructure:"MYSQL_DATABASE" json:"database"`
	User     string `mapstructure:"MYSQL_USER_NAME" json:"user_name"`
	Passwd   string `mapstructure:"MYSQL_PASSWD" json:"passwd"`
	Location string `mapstructure:"MYSQL_MAPPER_LOCATION" json:"mapper_location"`
}

func (c *MySQLConf) connString(passwd string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, c.User, passwd, c.Server, c.Database)
}

var Settings = &MySQLConf{
	Enable:   false,
	Name:     "default",
	Server:   "127.0.0.1:3306",
	Database: "mydb",
	User:     "keesh",
	Passwd:   "Cc",
}

var Databases []MySQLConf

func Init() {
	if !Settings.Enable {
		Log.Infof("MySQL-Disabled: val=%t", Settings.Enable)

		return
	}

	db, err := sql.Open("mysql", Settings.connString(Settings.Passwd))

	if err != nil {
		Log.Errorf("MySQL-failed: Error=%s, ConnectionString=%s", err.Error(), Settings.connString("****"))

		return
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		Log.Error(fmt.Sprintf("MySQL-failed: Error=%s, ConnectionString=%s", err.Error(), Settings.connString("****")))
	} else {
		var location = filepath.Join(utils.RootDir(), Settings.Location)

		if tpl, err := template.New("").Funcs(funcMaps()).ParseGlob(location); err != nil {
			panic(err)
		} else {
			Client = MySQLClient{
				DB:     db,
				render: tpl,
				conf:   Settings,
			}
		}

		if version, err := Client.Version(); err == nil {
			Log.Info(fmt.Sprintf("MySQL-Succeed: ServerVersion=%s, ConnString=%s", version, Settings.connString("****")))
		}
	}

	// 构建多数据源
	SetupDataSource()
}

func SetupDataSource() {
	if len(Databases) == 0 {
		return
	}

	DataSouces = make(map[string]*MySQLClient, len(Databases))

	for _, m := range Databases {
		db, err := sql.Open("mysql", m.connString(m.Passwd))

		if err != nil {
			Log.Errorf("MySQL-failed: Name=%s, Error=%s, ConnectionString=%s", m.Name, err.Error(), m.connString("****"))

			continue
		}

		// See "Important settings" section.
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		if err := db.Ping(); err != nil {
			Log.Error(fmt.Sprintf("MySQL-failed: Name=%s, Error=%s, ConnectionString=%s", m.Name, err.Error(), m.connString("****")))

			continue
		}

		var location = filepath.Join(utils.RootDir(), m.Location)

		if tpl, err := template.New("").Funcs(funcMaps()).ParseGlob(location); err != nil {
			Log.Error(fmt.Sprintf("MySQL-failed: Name=%s, Error=%s, ConnectionString=%s", m.Name, err.Error(), m.connString("****")))

			continue
		} else {
			client := MySQLClient{
				DB:     db,
				render: tpl,
				conf:   &m,
			}

			if version, err := client.Version(); err == nil {
				Log.Info(fmt.Sprintf("MySQL-Succeed: Name=%s, ServerVersion=%s, ConnString=%s", m.Name, version, m.connString("****")))
			}

			DataSouces[m.Name] = &client
		}
	}
}
