package client

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io/ioutil"
)

type Client struct {
    HttpClient *http.Client
}

func NewClient() *Client {
    return &Client{
        HttpClient: &http.Client{},
    }
}

func (c *Client) MakeRequest(method, url string, payload interface{}) ([]byte, error) {
    var reqBody []byte
    var err error

    if payload != nil {
        reqBody, err = json.Marshal(payload)
        if err != nil {
            return nil, err
        }
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}
