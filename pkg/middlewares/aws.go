package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

const (
	signatureVersion       = "2"
	signatureMethod        = "HmacSHA256"
	timeFormatZ            = "2006-01-02T15:04:05Z"
	timeFormat2            = "2006-01-02T15:04:05"
	ExpireSecondsFieldKey  = "ExpireSeconds"
	AWSAccessKeyIdFieldKey = "AWSAccessKeyId"
)

func SignV2Request(req *http.Request, accessKeyID string, secretAccessKey string, expireSeconds int) {
	preSignedUrl, err1 := PreSignedUrlV2(accessKeyID, secretAccessKey, expireSeconds, req.Method, req.URL, nil)

	if err1 != nil {
		panic(err1)
	}

	req.URL = preSignedUrl
}

func PreSignedUrlV2(accessKeyID string, secretAccessKey string, expireSeconds int, method string, reqUrl *url.URL, queryParams *map[string]interface{}) (*url.URL, error) {
	var Query = reqUrl.Query()

	if queryParams != nil {
		for key, val := range *queryParams {
			Query.Set(key, utils.ToString(val))
		}
	}

	// Set new query parameters
	Query.Set(AWSAccessKeyIdFieldKey, accessKeyID)
	Query.Set("SignatureVersion", signatureVersion)
	Query.Set("SignatureMethod", signatureMethod)

	if !Query.Has("Timestamp") {
		Query.Set("Timestamp", time.Now().UTC().Format(timeFormatZ))
	}

	if expireSeconds > 0 {
		Query.Set(ExpireSecondsFieldKey, utils.ToString(expireSeconds))
	}

	// in case this is a retry, ensure no signature present
	Query.Del("Signature")

	host := reqUrl.Host
	path := reqUrl.Path

	if path == "" {
		path = "/"
	} else if strings.Contains(path, "../") {
		return nil, errors.New("ILLEGAL_SIGNATURE")
	}

	// obtain all of the query keys and sort them
	queryKeys := make([]string, 0, len(Query))

	for key := range Query {
		queryKeys = append(queryKeys, key)
	}

	// sort keys
	sort.Strings(queryKeys)

	// build URL-encoded query keys and values
	queryKeysAndValues := make([]string, len(queryKeys))

	for i, key := range queryKeys {
		k := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		v := strings.Replace(url.QueryEscape(Query.Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = k + "=" + v
	}

	// join into one query string
	query := strings.Join(queryKeysAndValues, "&")

	// build the canonical string for the V2 signature
	stringToSign := strings.Join([]string{
		strings.ToUpper(method),
		host,
		path,
		query,
	}, "\n")

	hash := hmac.New(sha256.New, []byte(secretAccessKey))
	hash.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	Query.Set("Signature", signature)

	newUrl, err1 := url.Parse(fmt.Sprintf(`%s://%s%s?%s`, reqUrl.Scheme, reqUrl.Host, reqUrl.Path, Query.Encode()))

	if err1 != nil {
		panic(err1)
	}

	return newUrl, nil
}

func CheckSign(accessKeyID string, secretAccessKey string, method string, reqUrl *url.URL) (bool, error) {
	var accessQuery = reqUrl.Query()

	if !accessQuery.Has("Signature") {
		return false, errors.New("ILLEGAL_SIGNATURE")
	}

	expireSeconds, _ := strconv.Atoi(accessQuery.Get(ExpireSecondsFieldKey))

	preSignedUrl, err1 := PreSignedUrlV2(accessKeyID, secretAccessKey, expireSeconds, method, reqUrl, nil)

	if err1 != nil {
		return false, err1
	}

	var newQuery url.Values = preSignedUrl.Query()

	// 签名不匹配, 签名无效
	if accessQuery.Get("Signature") != newQuery.Get("Signature") {
		return false, errors.New("INVALID_SIGNATURE")
	}

	if accessQuery.Has(ExpireSecondsFieldKey) {
		expireSeconds, err0 := strconv.Atoi(accessQuery.Get(ExpireSecondsFieldKey))

		// 过期时间格式不对
		if err0 != nil {
			return false, errors.New("ILLEGAL_EXPIRES_ECONDS")
		}

		// 无过期时间
		if expireSeconds <= 0 {
			return true, nil
		}

		var signTime time.Time
		var err error
		var timestampString string = accessQuery.Get("Timestamp")

		if strings.HasSuffix(timestampString, "Z") {
			signTime, err = time.Parse(timeFormatZ, timestampString)
		} else {
			signTime, err = time.Parse(timeFormat2, timestampString)
		}

		// 时间格式不对
		if err != nil {
			return false, errors.New("ILLEGAL_INVALID_TIMESTAMP")
		}

		// 过期时间在当前时间之后, 签名有效
		if signTime.Add(time.Duration(expireSeconds) * time.Second).After(time.Now()) {
			return true, nil
		}

		// 签名已过期
		return false, errors.New("SIGNATURE_EXPIRED")
	}

	return true, nil
}
