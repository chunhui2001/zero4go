package gsql

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"text/template"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
	"github.com/chunhui2001/zero4go/pkg/utils"
)

var ReGTrim = regexp.MustCompile(`,\s*__TRIM__\(,\)`)

// SqlBindContext 模拟 Rust 的 Arc<Mutex<Vec<Value>>>
type SqlBindContext struct {
	binds []interface{}
	mu    sync.Mutex
}

func NewSqlBindContext() *SqlBindContext {
	return &SqlBindContext{}
}

// AddBind 添加一个绑定值
func (ctx *SqlBindContext) AddBind(val interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.binds = append(ctx.binds, val)
}

// TakeBinds 获取全部绑定值（取出后清空）
func (ctx *SqlBindContext) TakeBinds() []interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	result := ctx.binds
	ctx.binds = nil

	return result
}

func funcMaps() template.FuncMap {
	return template.FuncMap{
		"sql_bind": SqlBind,
		"TrimComma": func(s string) string {
			return "__TRIM__(" + strings.TrimSpace(s) + ")"
		},
	}
}

// SqlBind
// AND name = {{ sql_bind .Name .BindCtx }}
// AND id   IN ({{ sql_bind_in .Ids .BindCtx }})
func SqlBind(val any, ctx *SqlBindContext) string {
	if ctx == nil {
		panic("sql_bind: ctx is nil")
	}

	ctx.AddBind(val)

	return "?"
}

// RenderSQL 1️⃣ 渲染 SQL + 获取绑定值
func RenderSQL(tmpl *template.Template, tplName string, params map[string]interface{}) (string, []any, error) {
	bindCtx := NewSqlBindContext()

	params["CTX"] = bindCtx

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, tplName, params)

	if err != nil {
		Log.Errorf("RenderSQL: Error=%s", err.Error())

		return "", nil, err
	}

	// 取出绑定值
	binds := bindCtx.TakeBinds()

	var out = buf.String()

	var sqlStatement = trimEndSymbol(utils.NormalizeSpace(DebugSQLWithBinds(out, binds)))

	Log.Debugf("sqlStatement: sql=%s", sqlStatement)

	return trimEndSymbol(out), binds, nil
}

// DebugSQLWithBinds 将 SQL 中的 '?' 占位符替换成绑定值的可读形式
func DebugSQLWithBinds(sql string, binds []interface{}) string {
	var sb strings.Builder
	bindIdx := 0

	for _, c := range sql {
		if c == '?' {
			if bindIdx < len(binds) {
				val := binds[bindIdx]
				bindIdx++

				switch v := val.(type) {
				case string:
					// 字符串加单引号，并替换内部单引号为双单引号
					escaped := strings.ReplaceAll(v, "'", "''")
					sb.WriteString("'" + escaped + "'")
				case fmt.Stringer:
					sb.WriteString(v.String())
				default:
					sb.WriteString(fmt.Sprintf("%v", v))
				}
			} else {
				sb.WriteRune('?')
			}
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

func trimEndSymbol(out string) string {
	return ReGTrim.ReplaceAllString(out, "")
}
