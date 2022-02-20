package mongoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	ActionFindOne = "findOne"
	ActionFind    = "find"
)

var (
	ErrNoDocuments = errors.New("mongoapi: no documents found in result")
)

type Client struct {
	url        string
	key        string
	httpclient *http.Client
}

type Database struct {
	client   *Client
	cluster  string
	database string
}

type Collection struct {
	name     string
	database *Database
}

func New(url, key string) *Client {
	cl := http.DefaultClient

	cl.Timeout = 90 * time.Second

	// TODO: custom http client

	return &Client{
		url:        url,
		key:        key,
		httpclient: cl,
	}
}

func (c *Client) Database(cluster, database string) *Database {
	return &Database{
		client:   c,
		cluster:  cluster,
		database: database,
	}
}

func (d *Database) Collection(name string) *Collection {
	return &Collection{
		name:     name,
		database: d,
	}
}

type ResponseError struct {
	Message string `json:"error"`
	Code    string `json:"error_code"`
	Link    string `json:"link"`
}

func (d *Database) do(ctx context.Context, action string, body bson.M, dest interface{}) error {
	body["dataSource"] = d.cluster
	body["database"] = d.database

	bodyStr, err := bson.MarshalExtJSON(body, true, false)
	if err != nil {
		return fmt.Errorf("mongoapi: failed to marshal bson request body: %w", err)
	}

	fullURL := fmt.Sprintf("%s/action/%s", d.client.url, action)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(string(bodyStr)))
	if err != nil {
		return fmt.Errorf("mongoapi: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Request-Headers", "*")
	req.Header["api-key"] = []string{d.client.key} // preserves casing

	res, err := d.client.httpclient.Do(req)
	if err != nil {
		return fmt.Errorf("mongoapi: request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("mongoapi: response body read error: %w", err)
	}

	if res.StatusCode == http.StatusOK {
		if len(resBodyBytes) > 0 {
			err = json.Unmarshal(resBodyBytes, &dest)
			if err != nil {
				return fmt.Errorf("mongoapi: response parsing failed: %w", err)
			}
			return nil
		} else {
			return fmt.Errorf("mongoapi: expected response body but none was found")
		}
	}

	errRes := ResponseError{}
	err = json.Unmarshal(resBodyBytes, &errRes)
	if err != nil {
		return fmt.Errorf("mongoapi: response parsing failed: %w", err)
	}

	return fmt.Errorf("mongoapi: request failed: %s %s", errRes.Message, errRes.Link)
}
