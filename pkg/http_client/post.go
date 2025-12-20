package http_client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"moul.io/http2curl"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

func HttpPost(reqUrl string, contentType string, data []byte) ([]byte, error) {

	myHttpClient := &http.Client{
		Transport: defaultTransport(),
		Timeout:   time.Duration(35) * time.Second,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}

	var req *http.Request

	req, _ = http.NewRequest("POST", reqUrl, bytes.NewBuffer(data))

	resp, err := myHttpClient.Do(req)

	command, _ := http2curl.GetCurlCommand(req)
	commandCurl := command.String()

	if Settings.HttpClientPrintCurl {
		Log.Debugf("commandCurl: curl=%s", commandCurl)
	}

	if err != nil {
		Log.Errorf("HTTP POST error: %v", err)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		Log.Errorf("HTTP POST error: %v", err)

		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		Log.Errorf("HTTP %d: %s", resp.StatusCode, string(body))

		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
