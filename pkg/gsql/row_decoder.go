package gsql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

type RowDecoder[T any] struct {
	typ       reflect.Type
	isStruct  bool
	isScalar  bool
	isMap     bool
	singleCol bool
}

func NewRowDecoder[T any](cols []string) (*RowDecoder[T], error) {
	typ := utils.TypeOf[T]()

	d := &RowDecoder[T]{
		typ:       typ,
		isStruct:  utils.IsStruct(typ),
		isScalar:  utils.IsScalar(typ),
		isMap:     utils.IsMapStringAny(typ),
		singleCol: len(cols) == 1,
	}

	// 非法组合校验
	if d.isScalar && !d.singleCol {
		return nil, fmt.Errorf("scalar type %s requires single column", typ)
	}

	if d.isMap && !d.singleCol {
		// map[string]any 支持多列，这里允许
		return d, nil
	}

	if !d.isStruct && !d.isScalar && !d.isMap {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}

	return d, nil
}

func (d *RowDecoder[T]) Decode(row *sql.Rows, cols []string) (T, error) {
	var zero T

	switch {
	case d.isScalar:
		var v T
		if err := row.Scan(&v); err != nil {
			return zero, err
		}

		return v, nil

	case d.isMap:
		return d.decodeMap(row, cols)

	case d.isStruct:
		ptr, err := d.mapColumns(row, cols)

		if err != nil {
			return zero, err
		}

		return *ptr, nil
	}

	return zero, fmt.Errorf("unreachable")
}

func (d *RowDecoder[T]) decodeMap(rows *sql.Rows, cols []string) (T, error) {
	values := make([]any, len(cols))
	scanArgs := make([]any, len(cols))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	if err := rows.Scan(scanArgs...); err != nil {
		var zero T
		return zero, err
	}

	m := make(map[string]any, len(cols))

	for i, col := range cols {
		v := values[i]
		// 处理 []byte → string（MySQL 常见）
		if b, ok := v.([]byte); ok {
			m[col] = string(b)
		} else {
			m[col] = v
		}
	}

	return any(m).(T), nil
}

func (d *RowDecoder[T]) mapColumns(row *sql.Rows, cols []string) (*T, error) {
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

	if err := row.Scan(values...); err != nil {
		return nil, err
	}

	return &result, nil
}
