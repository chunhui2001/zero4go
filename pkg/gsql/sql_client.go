package gsql

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
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
		Log.Errorf("MySQL-Insert-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Insert-Error: Error=%s", err.Error())

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
		Log.Errorf("MySQL-Update-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Update-Error: Error=%s", err.Error())

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
		Log.Errorf("MySQL-Delete-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("MySQL-Delete-Error: Error=%s", err.Error())

		return -1, err
	}

	return result.RowsAffected()
}

func (s *MySQLClient) SelectSingleRow(tplName string, params map[string]any) (map[string]any, error) {
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 4️⃣ 获取列名
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-QuerySingleRow-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	columns, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-QuerySingleRow-Error: Error=%s", err.Error())

		return nil, err
	}

	values := make([]interface{}, len(columns))

	for i := range values {
		var tmp interface{}
		values[i] = &tmp
	}

	if err := rows.Scan(values...); err != nil {
		Log.Errorf("MySQL-QuerySingleRow-Error: Error=%s", err.Error())

		return nil, err
	}

	rowMap := make(map[string]interface{}, len(columns))

	for i, col := range columns {
		raw := *(values[i].(*interface{}))

		switch v := raw.(type) {
		case []byte:
			// MySQL 字符串列返回 []byte，需要转成 string
			rowMap[col] = string(v)
		default:
			rowMap[col] = v
		}
	}

	return rowMap, nil
}

func SelectRow[T any](s *MySQLClient, tplName string, params map[string]any) (*T, error) {
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 4️⃣ 获取列名
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-SelectRow-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	cols, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-SelectRow-Error: Error=%s", err.Error())

		return nil, err
	}

	// 创建 T 实例
	result, err := mapColumns[T](rows, cols)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MySQLClient) SelectMultipleRows(tplName string, params map[string]any) ([]map[string]any, error) {
	// 1️⃣ 渲染 SQL + 获取绑定值
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 2️⃣ 执行查询
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-SelectMultipleRows-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	// 3️⃣ 获取列名
	columns, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-SelectMultipleRows-Error: Error=%s", err.Error())

		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// 4️⃣ 创建占位切片
		values := make([]interface{}, len(columns))

		for i := range values {
			var tmp interface{}
			values[i] = &tmp
		}

		// 5️⃣ Scan
		if err := rows.Scan(values...); err != nil {
			Log.Errorf("MySQL-SelectMultipleRows-Error: Error=%s", err.Error())

			return nil, err
		}

		// 6️⃣ 构造单行 map
		rowMap := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			raw := *(values[i].(*interface{}))
			switch v := raw.(type) {
			case []byte:
				rowMap[col] = string(v) // 转字符串
			default:
				rowMap[col] = v
			}
		}

		results = append(results, rowMap)
	}

	// 检查遍历过程中是否出错
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func SelectRows[T any](s *MySQLClient, tplName string, params map[string]any) ([]*T, error) {
	// 1️⃣ 渲染 SQL + 获取绑定值
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 2️⃣ 执行查询
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("MySQL-SelectRows-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	// 3️⃣ 获取列名
	cols, err := rows.Columns()

	if err != nil {
		Log.Errorf("MySQL-SelectRows-Error: Error=%s", err.Error())

		return nil, err
	}

	var results []*T

	for rows.Next() {
		var result, err = mapColumns[T](rows, cols)

		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func mapColumns[T any](rows *sql.Rows, cols []string) (*T, error) {
	var result T

	val := reflect.ValueOf(&result).Elem()
	typ := val.Type()

	// 字段名 → index 映射（支持 db tag）
	fieldMap := make(map[string]int)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		name := f.Tag.Get("db")

		if name == "" {
			name = strings.ToLower(f.Name)
		}

		fieldMap[name] = i
	}

	// scan 容器
	values := make([]any, len(cols))

	for i, col := range cols {
		if idx, ok := fieldMap[col]; ok {
			values[i] = val.Field(idx).Addr().Interface()
		} else {
			var dummy any
			values[i] = &dummy
		}
	}

	if err := rows.Scan(values...); err != nil {
		Log.Errorf("Mysql-mapColumns-Error: Error=%s", err.Error())
		return nil, err
	}

	return &result, nil
}
