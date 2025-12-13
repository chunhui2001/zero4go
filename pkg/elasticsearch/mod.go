package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"

	. "github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/utils"
)

type ESConf struct {
	Enable      bool   `mapstructure:"ES_ICCO_ENABLE"`
	Servers     string `mapstructure:"ES_ICCO_SERVERS"`
	DslLocation string `mapstructure:"ES_ICCO_DSL_LOCATION"`
}

var Settings = &ESConf{
	Enable:      true,
	Servers:     "http://localhost:9200",
	DslLocation: "./META-INF/es_icco_dsl",
}

type EsClient struct {
	*elasticsearch.Client
}

var Client *EsClient

func Init() {
	if !Settings.Enable {
		return
	}

	retryBackoff := backoff.NewExponentialBackOff()

	cfg := elasticsearch.Config{
		MaxRetries:    5,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}

			return retryBackoff.NextBackOff()
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 5 * time.Second,
		},
		Addresses: strings.Split(Settings.Servers, ","),
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		Log.Errorf("ElasticSearch init Failed: Server=%s, Error=%s", Settings.Servers, err.Error())

		return
	}

	Client = &EsClient{
		es,
	}

	Client.Ping()
}

func (c *EsClient) Ping() {
	res, err := c.Info()

	if err != nil {
		Log.Errorf("ElasticSearch-Ping-Failed: server=%s, Error=%s", Settings.Servers, err.Error())

		return
	}

	defer res.Body.Close()

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		Log.Errorf("ElasticSearch-Ping-Failed: Error=%s", err.Error())
	} else {
		clusterName := r["cluster_name"].(string)
		serverInfo := r["version"].(map[string]interface{})
		serverVersion := serverInfo["number"]
		luceneVersion := serverInfo["lucene_version"]

		Log.Info(fmt.Sprintf(
			"Elastic-Init-Succeed: Servers=%s, ClusterName=%s, ServerVersion=%s, LuceneVersion=%s, ClientDriverVersion=%s",
			Settings.Servers, clusterName, serverVersion, luceneVersion, elasticsearch.Version),
		)
	}
}

// CatIndices 查询所有索引
func (c *EsClient) CatIndices() ([]map[string]interface{}, error) {

	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), c)

	if err != nil {
		Log.Errorf("Es-CatIndices-Failed: Error=%s", err.Error())

		return nil, err
	}

	defer res.Body.Close()
	var resMap []map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		Log.Errorf("Es-CatIndices-Failed: Error=%s", err.Error())

		return nil, err
	}

	return resMap, nil
}

// Save 新增
func (c *EsClient) Save(indexName string, dataMap map[string]interface{}) (string, error) {
	return c.SaveOrUpdate(indexName, "", dataMap)
}

// SaveOrUpdate 新增或更新
func (c *EsClient) SaveOrUpdate(indexName string, id string, dataMap map[string]interface{}) (string, error) {

	if dataMap == nil {
		return "", nil
	}

	if dataMap["@timestamp"] == nil {
		dataMap["@timestamp"] = utils.DateTimeUTCString()
	}

	_id := id

	if id == "" {
		_id = utils.Base64UUID()
	}

	// Instantiate a request object
	res, err := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: _id,
		Body:       strings.NewReader(utils.ToJsonString(dataMap)),
		Refresh:    "true",
	}.Do(context.Background(), c)

	if err != nil {
		Log.Errorf("Es-SaveOrUpdate-Failed: Error=%s", err.Error())

		return "", err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		Log.Errorf("Es-SaveOrUpdate-Failed: Error=%s", err.Error())

		return "", err
	}

	if resMap["error"] != nil {
		Log.Errorf("Es-SaveOrUpdate-Failed: Error=%s", utils.ToJsonString(resMap["error"]))

		return "", errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	return _id, nil
}

// Bulk 批量处理
func (c *EsClient) Bulk(indexName string, dataMap *[]map[string]interface{}) (bool, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        c,
		NumWorkers:    4,
		FlushBytes:    1024 * 1024, // bytes
		FlushInterval: 1 * time.Second,
	})

	if err != nil {
		Log.Errorf("Es-Bulk-Failed: Error=%s", err.Error())
		return false, err
	}

	var countSuccessful uint64

	for _, item := range *dataMap {

		err = bi.Add(context.Background(), c.getBulkIndexerItem(&item, &countSuccessful))

		if err != nil {
			panic(err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		panic(err)
	}

	biStatus := bi.Stats()

	if biStatus.NumFailed > 0 {
		return false, nil
	}

	return true, nil
}

func (c *EsClient) getBulkIndexerItem(item *map[string]interface{}, countSuccessful *uint64) esutil.BulkIndexerItem {

	if (*item)["_id"] == nil {
		(*item)["_id"] = utils.Base64UUID()
	}

	if (*item)["@timestamp"] == nil {
		(*item)["@timestamp"] = utils.DateTimeUTCString()
	}

	data, err := json.Marshal(item)

	if err != nil {
		panic(err)
	}

	return esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: (*item)["_id"].(string),
		Body:       bytes.NewReader(data),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			atomic.AddUint64(countSuccessful, 1)
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				Log.Errorf("Es-getBulkIndexerItem-Failed: Error=%s", err.Error())
			} else {
				Log.Errorf("Es-getBulkIndexerItem-Failed: ErrorType=%s, Error=%s", res.Error.Type, res.Error.Reason)
			}
		},
	}
}

func (c *EsClient) DoSearch(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {
	// Check for JSON errors
	// Default query is "{}" if JSON is invalid
	if !json.Valid([]byte(queryJsonString)) {
		Log.Errorf("Es-DoSearch-Failed: Error=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)

		return nil, 0, errors.New("not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := c.Search(
		c.Search.WithContext(context.Background()),
		c.Search.WithIndex(indexName),
		c.Search.WithBody(strings.NewReader(queryJsonString)),
		c.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		Log.Errorf("Es-DoSearch-Failed: queryJsonString=%s, Error=%s", queryJsonString, err.Error())

		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		Log.Errorf("Es-DoSearch-Failed: Error=%s", err.Error())

		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}

		Log.Errorf("Es-DoSearch-Failed: Error=%s", utils.ToJsonString(resMap["error"]))

		return nil, 0, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil
	}

	hitsMap := resMap["hits"].(map[string]interface{})

	if hitsMap["hits"] == nil {
		return nil, 0, nil
	}

	var dataMap = hitsMap["hits"].([]interface{})
	total := hitsMap["total"].(map[string]interface{})["value"].(float64)

	var interfaceSlice []map[string]interface{}

	if total > 0 {
		for _, item := range dataMap {

			_map := item.(map[string]interface{})
			id := _map["_id"].(string)
			object := _map["_source"].(map[string]interface{})

			//Log.Infof(`Id=%s, len=%d`, id, len(object))

			object["id"] = id

			interfaceSlice = append(interfaceSlice, object)
		}
	}

	return interfaceSlice, int64(total), nil
}

func (c *EsClient) ConstructQuery(q string, size int) *strings.Reader {
	var queryJsonString = fmt.Sprintf(`{"query": { %s }, "size": %d}`, q, size)

	// Check for JSON errors
	// Default query is "{}" if JSON is invalid
	if !json.Valid([]byte(queryJsonString)) {
		Log.Warnf("constructQuery() ERROR: query string not valid: %s", queryJsonString)
		Log.Warnf("Using default match_all query")

		queryJsonString = "{}"
	}

	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(queryJsonString)

	// Instantiate a *strings.Reader object from string
	read := strings.NewReader(b.String())

	// Return a *strings.Reader object
	return read
}
