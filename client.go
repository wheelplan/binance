package market

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	key        string
	secret     string
	httpClient *http.Client
}

func NewClient(key, secret string) (c *Client) {
	client := &Client{
		key:        key,
		secret:     secret,
		httpClient: &http.Client{},
	}
	return client
}

func (c *Client) do(method, url string, auth bool, result interface{}) (resp *http.Response, err error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	if auth {

		if len(c.key) == 0 || len(c.secret) == 0 {
			err = errors.New("Private endpoints require you to set an API Key and API Secret")
			return
		}

		req.Header.Add("X-MBX-APIKEY", c.key)

		parms := req.URL.Query()

		timestamp := time.Now().Unix() * 1000
		parms.Set("timestamp", fmt.Sprintf("%d", timestamp))

		hashed := hmac.New(sha256.New, []byte(c.secret))
		_, err := hashed.Write([]byte(parms.Encode()))
		if err != nil {
			return nil, err
		}

		signature := hex.EncodeToString(hashed.Sum(nil))

		req.URL.RawQuery = parms.Encode() + "&signature=" + signature

	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = handleError(resp)
	if err != nil {
		return
	}

	if resp != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
	}

	return
}

func handleError(resp *http.Response) error {
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Bad response Status %s. Response Body: %s", resp.Status, string(body))
	}
	return nil
}

