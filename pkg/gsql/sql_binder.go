package gsql

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/utils"
	"github.com/flosch/pongo2/v6"
)

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

func RegisterSqlBindFilter() {
	pongo2.RegisterFilter("sql_bind", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		val := in.Interface()
		raw := param.Interface()

		bindCtx, ok := raw.(*SqlBindContext)

		if !ok || bindCtx == nil {
			panic("sql_bind: second argument must be *SqlBindContext")
		}

		bindCtx.AddBind(val)

		// 返回占位符
		return pongo2.AsValue("?"), nil
	})
}

// RenderSQL 1️⃣ 渲染 SQL + 获取绑定值
func RenderSQL(tplSet *pongo2.TemplateSet, tplName string, params map[string]interface{}) (string, []any, error) {
	// 获取模板
	// 2️⃣ 从模板文件加载模板
	tpl, err := tplSet.FromFile(tplName)

	if err != nil {
		Log.Errorf("RenderSQL: Error=%s", err.Error())

		return "", nil, err
	}

	bindCtx := NewSqlBindContext()

	params["bindCtx"] = bindCtx

	// 5️⃣ 执行模板
	out, err := tpl.Execute(params)

	if err != nil {
		Log.Errorf("RenderSQL: Error=%s", err.Error())

		return "", nil, err
	}

	// 取出绑定值
	binds := bindCtx.TakeBinds()

	var sqlStatement = utils.NormalizeSpace(DebugSQLWithBinds(out, binds))

	Log.Debugf("sqlStatement: sql=%s", sqlStatement)

	return out, binds, nil
}

// / 用于调试的函数：将 SQL 中的 `?` 占位符替换成绑定值的可读形式。

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
