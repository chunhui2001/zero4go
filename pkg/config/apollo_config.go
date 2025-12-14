package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/chunhui2001/zero4go/pkg/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/chunhui2001/zero4go/pkg/http_client"
)

func ReadApolloConfig(v *viper.Viper) {
	var apolloEnv = envName
	var apolloName = os.Getenv("APP_NAME")
	var apolloConfig = os.Getenv("APOLLO_CONFIG")
	var apolloSpace = os.Getenv("APOLLO_NAMESPACE")

	if len(apolloSpace) == 0 {
		apolloSpace = "application.properties,application.yaml"
	}

	log.Printf("Apollo config env: %s, %s, %s", apolloConfig, apolloName, apolloEnv)

	if len(apolloConfig) > 0 && len(apolloName) > 0 && len(apolloEnv) > 0 {
		for _, namespace := range strings.Split(apolloSpace, ",") {
			var array = strings.Split(namespace, ".")

			if len(array) == 2 {
				readNamespace(v, apolloConfig, apolloName, apolloEnv, namespace, array[1])
			} else {
				readNamespace(v, apolloConfig, apolloName, apolloEnv, namespace, "properties")
				readNamespace(v, apolloConfig, apolloName, apolloEnv, namespace, "yaml")
			}
		}
	}
}

func readNamespace(v *viper.Viper, apolloConfig string, apolloName string, apolloEnv string, namespace string, ext string) {
	var s1 = fmt.Sprintf("%s/configs/%s/%s/%s", apolloConfig, apolloName, apolloEnv, namespace)

	if res, err := http_client.SendRequest(s1); err != nil {

		log.Printf("Apollo Configuration Failed: Url=%s, Namespace=%s, Error=%s", s1, namespace, err.Error())
	} else {
		log.Printf("Apollo Loading Succeed: Url=%s", s1)

		if _map, _ := utils.ToMap(readResponse(res)); _map != nil {

			if ext == "properties" {
				readProperties(v, _map)
			} else if ext == "yaml" {
				readYaml(v, _map)
			}
		} else {

			log.Printf("Apollo Configuration Failed: Url=%s, Namespace=%s, Error=%s", s1, namespace, err.Error())
		}
	}
}

func readResponse(res *http.Response) []byte {
	data, _ := io.ReadAll(res.Body)

	var headerKey = os.Getenv("APOLLO_HEADER_KEY")
	var secretKey = os.Getenv("APOLLO_SECRET_KEY")

	if res.Header.Get(headerKey) == "true" {
		// 解密
		raw, err := base64.RawStdEncoding.DecodeString(secretKey)

		if err != nil {
			log.Printf("Apollo Configuration Failed: Error=%s", err.Error())

			return data
		}

		sk, err := x509.ParsePKCS8PrivateKey(raw)

		if err != nil {
			log.Printf("Apollo Configuration Failed: Error=%s", err.Error())

			return data
		}

		privKey := sk.(*rsa.PrivateKey)
		partLen := privKey.PublicKey.N.BitLen() / 8

		chunks := split(data, partLen)

		buffer := bytes.NewBufferString("")

		for _, chunk := range chunks {
			decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privKey, chunk)

			buffer.Write(decrypted)
		}

		return buffer.Bytes()
	}

	return data
}

func readProperties(v *viper.Viper, responseMap map[string]any) {
	config := responseMap["configurations"].(map[string]any)

	for key, val := range config {
		v.SetDefault(strings.TrimSpace(key), strings.TrimSpace(val.(string)))
	}
}

func readYaml(v *viper.Viper, responseMap map[string]any) {
	contents := responseMap["configurations"].(map[string]interface{})["content"].(string)

	var config map[string]any

	if err := yaml.Unmarshal([]byte(contents), &config); err != nil {
		log.Printf(`Loading-Yaml-File-Error: Contents=%s, Error=%s`, contents, err.Error())
	}

	for key, val := range NestedMap(config) {
		v.SetDefault(strings.TrimSpace(key), val)
	}
}

func NestedMap(input map[string]any) map[string]any {
	result := map[string]any{}

	// 正则匹配数组形式字段
	reArray := regexp.MustCompile(`^([^\[]+)\[(\d+)\]\.(.+)$`)

	for key, val := range input {
		if matches := reArray.FindStringSubmatch(key); len(matches) == 4 {
			listKey, indexStr, fieldPath := matches[1], matches[2], matches[3]
			index, _ := strconv.Atoi(indexStr)

			// 初始化数组容器
			if _, ok := result[listKey]; !ok {
				result[listKey] = map[int]map[string]any{}
			}

			list := result[listKey].(map[int]map[string]any)

			if _, ok := list[index]; !ok {
				list[index] = map[string]any{}
			}

			// 解析字段路径
			assignNestedMap(list[index], fieldPath, val)
		} else {
			// 普通嵌套 map 形式
			assignNestedMap(result, key, val)
		}
	}

	// 把 map[int] → []map 的转换进行处理
	for key, val := range result {
		if m, ok := val.(map[int]map[string]any); ok {
			length := len(m)
			list := make([]map[string]any, length)
			var indices []int

			for i := range m {
				indices = append(indices, i)
			}

			sort.Ints(indices)

			for i, idx := range indices {
				list[i] = m[idx]
			}

			result[key] = list
		}
	}

	// 输出为 YAML
	yamlBytes, err := yaml.Marshal(result)

	if err != nil {
		panic(err)
	}

	var body map[string]any

	if err := yaml.Unmarshal(yamlBytes, &body); err != nil {
		return body
	} else {
		return body
	}
}

// assignNestedMap 将路径如 "a.b.c" 分配到嵌套结构中
func assignNestedMap(m map[string]any, path string, val any) {
	parts := strings.Split(path, ".")
	curr := m

	for i, p := range parts {
		if i == len(parts)-1 {
			curr[p] = val
		} else {
			if _, ok := curr[p]; !ok {
				curr[p] = map[string]any{}
			}

			curr = curr[p].(map[string]any)
		}
	}
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)

	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}

	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}

	return chunks
}
