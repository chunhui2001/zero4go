package gsql

import (
	"context"
	"database/sql"
	"text/template"
	"time"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type MySQLClient struct {
	*sql.DB
	//render *pongo2.TemplateSet
	render *template.Template
}

func (s *MySQLClient) Version() (string, error) {
	var version string
	err2 := s.QueryRow("SELECT VERSION()").Scan(&version)

	return version, err2
}

func (s *MySQLClient) Insert(tplName string, params map[string]any) (int64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return 0, err
	}

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("MySQL-Insert-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Insert-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	return result.LastInsertId()
}

func (s *MySQLClient) Update(tplName string, params map[string]any) (int64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return 0, err
	}

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("MySQL-Update-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Update-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	return result.RowsAffected()
}

func (s *MySQLClient) Delete(tplName string, params map[string]any) (int64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return 0, err
	}

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("MySQL-Delete-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Delete-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return -1, err
	}

	return result.RowsAffected()
}

func SelectRow[T any](tplName string, params map[string]any) (*T, error) {
	// 1️⃣ 渲染 SQL + 获取绑定值
	sqlStr, binds, err := RenderSQL(Client.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 4️⃣ 获取列名
	rows, err := Client.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-SelectRow-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	cols, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-SelectRow-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return nil, err
	}

	decoder, err := NewRowDecoder[T](cols)

	if err != nil {
		return nil, err
	}

	// 创建 T 实例
	result, err := decoder.Decode(rows, cols)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SelectRows[T any](tplName string, params map[string]any) ([]T, error) {
	// 1️⃣ 渲染 SQL + 获取绑定值
	sqlStr, binds, err := RenderSQL(Client.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 2️⃣ 执行查询
	rows, err := Client.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-SelectRows-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	// 3️⃣ 获取列名
	cols, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-SelectRows-Error: tplName=%s, sqlStr=%s, Error=%s", tplName, sqlStr, err.Error())

		return nil, err
	}

	decoder, err := NewRowDecoder[T](cols)

	if err != nil {
		return nil, err
	}

	var results = make([]T, 0)

	for rows.Next() {
		v, err := decoder.Decode(rows, cols)

		if err != nil {
			return nil, err
		}

		results = append(results, v)
	}

	return results, nil
}
