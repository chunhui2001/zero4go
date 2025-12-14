package gsql

import (
	"context"
	"database/sql"
	"text/template"
	"time"

	. "github.com/chunhui2001/zero4go/pkg/logs"
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

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("Mysql-Insert-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("Mysql-Insert-Error: Error=%s", err.Error())

		return -1, err
	}

	return result.LastInsertId()
}

func (s *MySQLClient) Update(tplName string, params map[string]any) (int64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("Mysql-Update-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("Mysql-Update-Error: Error=%s", err.Error())

		return -1, err
	}

	return result.RowsAffected()
}

func (s *MySQLClient) Delete(tplName string, params map[string]any) (int64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	// prepare the statement
	stmt, err := s.Prepare(sqlStr)

	if err != nil {
		Log.Errorf("Mysql-Delete-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return -1, err
	}

	defer stmt.Close()

	// format all args at once
	result, err := stmt.ExecContext(ctx, binds...)

	if err != nil {
		Log.Errorf("Mysql-Delete-Error: Error=%s", err.Error())

		return -1, err
	}

	return result.RowsAffected()
}

// QueryRowSQL 执行 SQL 模板查询，返回单行结果 map
func (s *MySQLClient) QuerySingleRow(tplName string, params map[string]any) (map[string]any, error) {
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 4️⃣ 获取列名
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("Mysql-QuerySingleRow-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	columns, err := rows.Columns()

	if err != nil {
		Log.Errorf("Mysql-QuerySingleRow-Error: Error=%s", err.Error())

		return nil, err
	}

	values := make([]interface{}, len(columns))

	for i := range values {
		var tmp interface{}
		values[i] = &tmp
	}

	if err := rows.Scan(values...); err != nil {
		Log.Errorf("Mysql-QuerySingleRow-Error: Error=%s", err.Error())

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

func (s *MySQLClient) QueryMultipleRows(tplName string, params map[string]any) ([]map[string]any, error) {
	// 1️⃣ 渲染 SQL + 获取绑定值
	sqlStr, binds, err := RenderSQL(s.render, tplName, params)

	if err != nil {
		return nil, err
	}

	// 2️⃣ 执行查询
	rows, err := s.Query(sqlStr, binds...)

	if err != nil {
		Log.Errorf("Mysql-QueryMultipleRows-Error: sqlStr=%s, Error=%s", sqlStr, err.Error())

		return nil, err
	}

	defer rows.Close()

	// 3️⃣ 获取列名
	columns, err := rows.Columns()

	if err != nil {
		Log.Errorf("Mysql-QueryMultipleRows-Error: Error=%s", err.Error())

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
			Log.Errorf("Mysql-QueryMultipleRows-Error: Error=%s", err.Error())

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
