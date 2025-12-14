package gzook

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/go-zookeeper/zk"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

type ZookConf struct {
	Enabled bool   `mapstructure:"ZOOKEEPER_ENABLED"`
	Servers string `mapstructure:"ZOOKEEPER_SERVERS"`
}

var Settings = &ZookConf{
	Enabled: true,
	Servers: "127.0.0.1:2181",
}

var Client *ZooClient

type ZooClient struct {
	*zk.Conn
}

func Init() {
	if !Settings.Enabled {
		return
	}

	var timeOut = time.Duration(5) * time.Second
	conn, _, err := zk.Connect(strings.Split(Settings.Servers, ","), timeOut)

	if err != nil {
		Log.Errorf(`ZooKeeper-Init-Failed: ConnectTimeout=%s, Servers=%s, Error=%s`, timeOut, Settings.Servers, err.Error())
		return
	}

	Client = &ZooClient{
		conn,
	}

	Log.Infof(`ZooKeeper-Connect-Succeed: ConnectTimeout=%s, Servers=%s, SessionId=%d`, timeOut, Settings.Servers, conn.SessionID())

	Client.Info()
}

func (z *ZooClient) Info() {
	var addrs = strings.Split(Settings.Servers, ",")

	var lines = make([]string, 0)

	for _, addr := range addrs {
		var status = z.Srvr(addr)

		re := regexp.MustCompile(`(\d+\.\d+\.\d+)-([\w\W]*?)(Mode: ((follower)|(leader)))`)
		matches := re.FindStringSubmatch(status)
		var version string
		var mode string

		if len(matches) > 4 {
			version = matches[1]
			mode = matches[4]
		}

		lines = append(lines, fmt.Sprintf(`ZooKeeper-Info: addr=%s, version=%s, Mode=%+v`, addr, version, mode))
	}

	for _, line := range lines {
		Log.Info(line)
	}
}

func (z *ZooClient) TryLock(path string, data string, timeout time.Duration) {

	thePath := path

	if !strings.HasPrefix(path, "/") {
		thePath = "/" + path
	}

	for !z.tryLock(thePath, data) {
		logger.Debugf(`ZooKeeper-TryLock-Failed-Retry: LockTimeout=%s, LockPath=%s, Data=%s`, timeout, thePath, data)
	}
}

// flags 有4种取值：
// 0:永久，除非手动删除
// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
// zk.FlagSequence  = 2:会自动在节点后面添加序号
// 3:Ephemeral和Sequence，即，短暂且自动添加序号
func (z *ZooClient) tryLock(path string, data string) bool {
	_, err := z.Create(path, []byte(data), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	if err == nil {
		Log.Infof(`ZooKeeper-TryLock-Succeed: LockPath=%s, Data=%s`, path, data)

		return true
	}

	_, _, watch, err := z.ExistsW(path)

	if err != nil {
		Log.Errorf(`ZooKeeper-Watch-Error: LockPath=%s, ErrorMessage=%s`, path, err.Error())
	}

	select {
	case event := <-watch:
		if event.Type == zk.EventNodeDeleted {
			_, err := z.Create(path, []byte(data), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err == nil {
				Log.Warnf(`ZooKeeper-TryLock-Succeed: LockPath=%s, Data=%s`, path, data)

				return true
			}
		}
	// time out
	case <-time.After(5 * time.Second):
		return false
	}

	return false
}

func (z *ZooClient) Read(path string) (string, error) {
	data, _, err := z.Get(path)

	return string(data), err
}

// Modify 删改与增不同在于其函数中的version参数,其中version是用于 CAS支持, 可以通过此种方式保证原子性
func (z *ZooClient) Modify(path string, data string) error {
	newData := []byte(data)

	_, sate, _ := z.Get(path)
	_, err := z.Set(path, newData, sate.Version)

	return err
}

func (z *ZooClient) Del(path string) error {
	_, sate, _ := z.Get(path)

	err := z.Delete(path, sate.Version)

	return err
}

func (z *ZooClient) Srvr(addr string) string {

	// 建立 TCP 连接
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		Log.Infof("连接失败: Addr=%s, Error=%s", addr, err.Error())

		return ""
	}

	defer conn.Close()

	// 四字节命令，例如 "srvr"
	cmd := "srvr"
	_, err = conn.Write([]byte(cmd))

	if err != nil {
		Log.Infof("发送命令失败: Addr=%s, Error=%s", addr, err.Error())

		return ""
	}

	// 读取返回数据
	scanner := bufio.NewScanner(conn)

	var sb strings.Builder

	for scanner.Scan() {
		sb.WriteString("\n")
		sb.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		Log.Infof("读取响应失败: Addr=%s, Error=%s", addr, err.Error())
	}

	return sb.String()
}
