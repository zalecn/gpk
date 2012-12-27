package gopack

import (
	"bytes"
	"encoding/json"
	"ericaro.net/gopack/protocol"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"errors"
)

func init() { // register this as a handler for file:/// url scheme
	protocol.RegisterClient("http", NewHttpClient)
}

type HttpClient struct {
	protocol.BaseClient
}

func NewHttpClient(name string, u url.URL, token *protocol.Token) (r protocol.Client, err error) {
	r = &HttpClient{
		BaseClient: *protocol.NewBaseClient(name, u, token),
	}
	return
}



func (c *HttpClient) Fetch(pid protocol.PID) (r io.ReadCloser, err error) {
	v := &url.Values{}
	pid.InParameter(v)
	//query url
	u := &url.URL{
		Path:     protocol.FETCH,
		RawQuery: v.Encode(),
	}
	remote := c.Path()
	resp, err := http.Get(remote.ResolveReference(u).String())
	if err != nil {
		return
	}
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		return nil, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (c *HttpClient) Push(pid protocol.PID, r io.Reader) (err error) {
	v := &url.Values{}
	pid.InParameter(v)
	//query url
	u := &url.URL{
		Path:     protocol.PUSH,
		RawQuery: v.Encode(),
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	var client http.Client
	remote := c.Path()
	req, err := http.NewRequest("POST", remote.ResolveReference(u).String(), buf)
	if err != nil {
		return
	}
	req.ContentLength = int64(buf.Len()) // fuck I can't do that, I need to compute the length first
	resp, err := client.Do(req)
	if resp.StatusCode <200 || resp.StatusCode >=300 {
		return errors.New(resp.Status)
	}
	return
}

func (c *HttpClient) Search(query string, start int) (result []protocol.PID) {
	v := url.Values{}
	v.Set("q", query)
	v.Set("start", strconv.Itoa(start))

	//query url
	u := &url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     protocol.SEARCH,
		RawQuery: v.Encode(),
	}
	remote := c.Path()
	resp, err := http.Get(remote.ResolveReference(u).String())
	if err != nil {
		return result
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return result
}
