package shttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/slib/stime"
)

const (
	FormatJson = "json"
	FormatText = "text"
)

type Client struct {
	client        *http.Client
	logger        *logrus.Logger                 // 如果为 nil，则不记录日志
	format        string                         // 响应报文格式
	error         error                          // 响应失败时的 JSON 结构。在响应成功和失败时 JSON 结构不一致时设置，不能是指针
	dumps         bool                           // 将 HTTP 报文打印到控制台
	baseUrl       string                         // Url 前缀
	logLength     int                            // 最大日志长度
	logEscape     bool                           // 是否转换日志中的特殊字符
	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *Client {
	cli := &http.Client{
		Timeout: 5 * time.Second,
	}

	client := &Client{
		client:        cli,
		logger:        nil,
		format:        FormatJson,
		error:         nil,
		dumps:         false,
		baseUrl:       "",
		logLength:     1024,
		logEscape:     false,
		requestBefore: nil,
		responseAfter: nil,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (cli *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinUrl(cli.baseUrl, path)
	}

	url2, err := AppendQuery(path, query)
	if err != nil {
		cli.logWarn("Append query failed, Url: %s, query: %s, error: %v.", url2, query.Encode(), err)
		return nil, fmt.Errorf("append query error [%v]", err)
	}

	req, err := http.NewRequest(method, url2, body)
	if err != nil {
		cli.logWarn("New request failed, Url: %s, error: %v.", url2, err)
		return nil, fmt.Errorf("new request error [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if cli.requestBefore != nil {
		cli.requestBefore(req)
	}

	return req, nil
}

func (cli *Client) doRequest(req *http.Request) (*http.Response, error) {
	cli.dumpRequest(req)

	res, err := cli.client.Do(req)
	if err != nil {
		cli.logWarn("Do request failed, Url: %s, error: %v.", req.URL, err)
		return nil, err
	}

	return res, nil
}

func (cli *Client) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	cli.dumpResponse(res)

	if cli.responseAfter != nil {
		err := cli.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error [%v]", err)
	}

	switch cli.format {
	case FormatText:
		return body, cli.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, cli.handleJsonResponse(res.StatusCode, body, out)
	}
}

func (cli *Client) handleJsonResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 == 2 {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					err2, ok2 := out.(error)
					if ok2 {
						return err2
					}
					return newResponseError(statusCode, string(body))
				}
			}
		}
		return nil
	} else {
		if cli.error != nil {
			ret := reflect.New(reflect.TypeOf(cli.error)).Interface()
			err := json.Unmarshal(body, ret)
			if err == nil {
				return ret.(error)
			}
		} else if out != nil {
			err := json.Unmarshal(body, out)
			if err == nil {
				err2, ok2 := out.(error)
				if ok2 {
					return err2
				}
			}
		}
		return newResponseError(statusCode, string(body))
	}
}

func (cli *Client) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 != 2 {
		return newResponseError(statusCode, string(body))
	}

	if out == nil {
		return nil
	}

	bs, ok := out.(*[]byte)
	if !ok {
		return fmt.Errorf("response receiver must *[]byte type")
	}
	*bs = body

	return nil
}

func (cli *Client) get(method string, path string, query url.Values, header http.Header, out interface{}) error {
	req, err := cli.newRequest(method, path, query, header, nil)
	if err != nil {
		return err
	}

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, nil, body, err, timer.Stops())

	return err
}

func (cli *Client) post(method string, path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := cli.newRequest(method, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, reqBody, resBody, err, timer.Stops())

	return err
}

func (cli *Client) Get(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodGet, path, query, header, out)
}

func (cli *Client) Post(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPost, path, query, header, in, out)
}

func (cli *Client) Put(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPut, path, query, header, in, out)
}

func (cli *Client) Delete(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodDelete, path, query, header, out)
}

func (cli *Client) GetBinary(path string, query url.Values, header http.Header) ([]byte, string, error) {
	req, err := cli.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return nil, "", err
	}

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	typ := res.Header.Get("Content-Type")

	cli.logRequest(req, res, nil, nil, err, timer.Stops())

	return dat, typ, err
}

func (cli *Client) PostJson(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.Post(path, query, header, in, out)
}

func (cli *Client) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()

	req, err := cli.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	reqBody, _ = url.QueryUnescape(reqBody)
	cli.logRequest(req, res, []byte(reqBody), resBody, err, timer.Stops())

	return err
}

func (cli *Client) PostFormData(path string, query url.Values, header http.Header, values map[string]string, files map[string]string, out interface{}) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
	}
	if len(files) > 0 {
		for key, file := range files {
			err := cli.writeFormFile(writer, key, file)
			if err != nil {
				return err
			}
		}
	}
	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	var cpy map[string]string
	var reqBody []byte
	if len(values) > 0 {
		cpy = values
		for key, file := range files {
			cpy[key] = file
		}
	} else {
		cpy = files
	}
	if len(cpy) == 0 {
		reqBody = []byte("{}")
	} else {
		reqBody, _ = json.Marshal(cpy)
	}
	cli.logRequest(req, res, reqBody, resBody, err, timer.Stops())

	return err
}

func (cli *Client) writeFormFile(writer *multipart.Writer, key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(key, file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) PostBinary(path string, query url.Values, header http.Header, mimeType string, in io.Reader, out interface{}) error {
	req, err := cli.newRequest(http.MethodPost, path, query, header, in)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mimeType)

	timer := stime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, nil, resBody, err, timer.Stops())

	return err
}
