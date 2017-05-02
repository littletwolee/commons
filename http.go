package commons

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	consHttp *Http
)

type Http struct{}

func GetHttp() *Http {
	if consHttp == nil {
		consHttp = &Http{}
	}
	return consHttp
}

// @Title HttpMethod
// @Description http client send http response
// @Parameters
//            data            *bytes.Reader                   data reader
//            method          string                          http method
//            url             string                          http url
//            headers         map[string]string               http headers
// @Returns statuscode:int json:[]byte err:error
func (h *Http) HttpMethod(data io.Reader, method, url string, headers map[string]string) (int, []byte, error) {
	request, err := http.NewRequest(method, url, data)
	code := -1
	if err != nil {
		return code, nil, err
	}
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return code, nil, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return code, nil, err
	}
	return resp.StatusCode, respBytes, nil
}

// @Title JsonData
// @Description http client send post
// @Parameters
//            data            map[string]interface{}          post json data
// @Returns reader:*bytes.Reader err:error
func (h *Http) JsonData(data map[string]interface{}) (*bytes.Reader, error) {
	var reader *bytes.Reader
	if data != nil && len(data) > 0 {
		bytesData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(bytesData)
	}
	return reader, nil
}
