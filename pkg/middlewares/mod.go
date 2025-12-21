package middlewares

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/utils"
)

type AccessClients struct {
	ClientName      string   `yaml:"clientName"`
	AccessKeyID     string   `yaml:"accessKeyID"`
	SecretAccessKey string   `yaml:"secretAccessKey"`
	Enabled         bool     `yaml:"enabled"`
	Scope           []string `yaml:"scope"`
}

var (
	accessClientsMap = make(map[string]*AccessClients)
)

func Init() {
	if v := config.GetConfig("Access-Clients"); v != nil {
		for key, val := range v.(map[string]interface{}) {
			if vv, err := utils.ToStructByTag[AccessClients](val, "yaml"); err == nil {
				accessClientsMap[key] = &vv
			}
		}
	}
}

func RequestURL(req *http.Request) *url.URL {

	scheme := "http"

	if req.TLS != nil {
		scheme = "https"
	}

	var rawQuery = req.URL.Query().Encode()

	if rawQuery != "" {
		newUrl, err1 := url.Parse(fmt.Sprintf(`%s://%s%s?%s`, scheme, req.Host, req.URL.Path, rawQuery))

		if err1 != nil {
			panic(err1)
		}

		return newUrl
	}

	newUrl, err1 := url.Parse(fmt.Sprintf(`%s://%s%s`, scheme, req.Host, req.URL.Path))

	if err1 != nil {
		panic(err1)
	}

	return newUrl
}
