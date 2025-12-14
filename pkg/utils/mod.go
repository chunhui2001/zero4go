package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

const TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func OutboundIP() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func NowTimestamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func DateTimeUTCString() string {
	return time.Now().UTC().Format(TimeStampFormat)
}

func DateOffsets() (string, string) {

	now := time.Now()
	_, offset := now.Zone()

	if offset >= 0 {
		return os.Getenv("TZ"), fmt.Sprintf("+%02d:00", offset/3600)
	}

	return os.Getenv("TZ"), fmt.Sprintf("-%02d:00", offset/3600)
}

func RootDir() string {
	var appRoot = os.Getenv("WORK_DIR")

	if appRoot != "" {
		return appRoot
	}

	dir, _ := os.Getwd()

	return dir
}

func TempDir() string {
	return os.TempDir()
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	return true, nil
}

func ToString(s any) string {

	switch s.(type) {
	case float64, float32:
		return fmt.Sprint(s)
	case string:
		return fmt.Sprintf("%s", s)
	case bool:
		return fmt.Sprintf("%t", s)
	case byte:
		return fmt.Sprintf("%x", s)
	case []uint8:
		return string(s.([]byte))
	default:
		return fmt.Sprintf("%d", s)
	}
}

func ToJsonString(v any) string {
	if v == nil {
		return ""
	}

	if reflect.TypeOf(v).String() == "string" {
		return v.(string)
	}

	b, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func StrToInt(str string) int {
	if str == "" {
		return 0
	}

	intVar, err := strconv.Atoi(str)

	if err != nil {
		log.Printf("StrToInt: value=%s, Error=%s", str, err.Error())

		panic(err)
	}

	return intVar
}

func NormalizeSpace(s string) string {
	// 1. 去掉所有换行符
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ") // Windows 换行

	// 2. 将一个或多个空格替换为单个空格
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// 3. 去掉首尾空格
	return strings.TrimSpace(s)
}

func ToDateTimeUTCString(tm time.Time) string {
	return tm.Format(TimeStampFormat)
}

func OfMap(kv ...any) map[string]any {
	if kv == nil {
		return make(map[string]any)
	}

	if len(kv)%2 != 0 {
		panic(errors.New("Invalid map size: currentSize=" + ToString(len(kv))))
	}

	m := make(map[string]any)

	for i := 0; i < len(kv); i++ {
		k := ToString(kv[i])
		m[k] = kv[i+1]

		i++
	}

	return m
}

func ToMap(buf []byte) (map[string]any, error) {
	var m map[string]any

	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func MapsToInterfaces(in []map[string]interface{}) []interface{} {
	out := make([]interface{}, len(in))

	for i, v := range in {
		out[i] = v
	}

	return out
}

func ToSlice[T any](input []map[string]interface{}) ([]T, error) {

	return ToSliceAny[T](MapsToInterfaces(input))
}

func ToSliceAny[T any](input []interface{}) ([]T, error) {
	result := make([]T, 0, len(input))

	for _, item := range input {
		s, err := ToStruct[T](item)

		if err != nil {
			return nil, err
		}

		result = append(result, s)
	}

	return result, nil
}

func ToStruct[T any](input interface{}) (T, error) {

	return ToStructByTag[T](input, "json")
}

func ToStructByTag[T any](input interface{}, tagName string) (T, error) {
	var result T

	config := &mapstructure.DecoderConfig{
		Result:           &result,
		TagName:          tagName,
		WeaklyTypedInput: true, // 支持弱类型转换
		DecodeHook: func(f reflect.Type, t reflect.Type, data any) (any, error) {
			if f.Kind() == reflect.String && t == reflect.TypeOf(decimal.Decimal{}) {
				s := data.(string)
				d, err := decimal.NewFromString(s)

				if err != nil {
					return nil, err
				}

				return d, nil
			}

			return data, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)

	if err != nil {
		return result, err
	}

	err = decoder.Decode(input)

	return result, err
}

func ToBase64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func FromBase64String(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Base64UUID() string {

	b := make([]byte, 16)

	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}

	b[6] &= 0x0f /* clear the 4 most significant bits for the version  */
	b[6] |= 0x40 /* set the version to 0100 / 0x40 */

	/* Set the variant:
	 * The high field of th clock sequence multiplexed with the variant.
	 * We set only the MSB of the variant*/
	b[8] &= 0x3f /* clear the 2 most significant bits */
	b[8] |= 0x80 /* set the variant (MSB is set)*/

	return base64.RawURLEncoding.EncodeToString(b)
}

func ReadFile(filepath string) []byte {
	content, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	return content
}

func ReadAllLines(filepath string) []string {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

func HumanFileSizeUint(size uint64) string {

	if size <= 0 {
		return "0"
	}

	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(float64(size)) / math.Log(1024)

	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)

	getSuffix := suffixes[int(math.Floor(base))]

	return strconv.FormatFloat(getSize, 'f', -1, 64) + "" + string(getSuffix)

}

func HumanFileSizeInt64(size int64) string {

	if size <= 0 {
		return "0"
	}

	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(float64(size)) / math.Log(1024)

	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)

	getSuffix := suffixes[int(math.Floor(base))]

	return strconv.FormatFloat(getSize, 'f', -1, 64) + "" + string(getSuffix)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow

	return
}
