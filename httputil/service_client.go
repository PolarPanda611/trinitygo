package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ServiceClient service client
type ServiceClient struct {
	Addr string
	Port int
}

// Request send request
func (s *ServiceClient) Request(method RequestMethod, path string, body []byte, header map[string]string) (int, interface{}, error) {
	url := fmt.Sprintf("http://%v:%v%v", s.Addr, s.Port, path)
	request, err := http.NewRequest(string(method), url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	//set default header
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(header) >= 1 {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	defer cancel()
	resp, err := client.Do(request.WithContext(ctx)) //发送请求
	if err != nil {
		return resp.StatusCode, nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	if len(respBytes) == 0 {
		return resp.StatusCode, nil, nil
	}
	var res interface{}
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		if err := json.Unmarshal(respBytes, &res); err != nil {
			return resp.StatusCode, nil, err
		}
		return resp.StatusCode, res, nil
	}
	return resp.StatusCode, nil, errors.New(string(respBytes))
}
