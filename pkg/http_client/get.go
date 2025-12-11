package http_client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"moul.io/http2curl"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

func HttpGet(reqUrl string) ([]byte, error) {

	myHttpClient := &http.Client{
		Transport: defaultTransport(),
		Timeout:   time.Duration(35) * time.Second,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}

	var req *http.Request

	req, _ = http.NewRequest("GET", reqUrl, nil)

	resp, err := myHttpClient.Do(req)

	command, _ := http2curl.GetCurlCommand(req)
	commandCurl := command.String()

	if Settings.HttpClientPrintCurl {
		Log.Debugf("commandCurl: curl=%s", commandCurl)
	}

	if err != nil {
		Log.Errorf("HTTP GET error: %v", err)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		Log.Infof("HTTP %d: %s", resp.StatusCode, string(body))

		return nil, errors.New(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
	}

	return body, nil
}
