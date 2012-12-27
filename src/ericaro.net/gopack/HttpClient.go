package gopack

import (
	"bytes"
	"encoding/json"
	. "ericaro.net/gopack/protocol"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func init() { // register this as a handler for file:/// url scheme
	h := func(name string, u url.URL) Client {
		f, _ := NewHttpClient(name, u)
		return f
	}
	RegisterClient("http", h)
}

func NewHttpClient(name string, u url.URL) (r *HttpClient, err error) {
	r = &HttpClient{
		remote: u,
		name:   name,
	}
	return

}

type HttpClient struct {
	remote url.URL
	name   string
}

func (r HttpClient) Name() string  { return r.name }
func (r HttpClient) Path() url.URL { return r.remote }
func (c *HttpClient) Fetch(pid PID) (r io.ReadCloser, err error) {
	v := &url.Values{}
	pid.InParameter(v)
	//query url
	u := &url.URL{
		Path:     FETCH,
		RawQuery: v.Encode(),
	}
	resp, err := http.Get(c.remote.ResolveReference(u).String())
	if err != nil {
		return
	}
	return resp.Body, nil
}

func (c *HttpClient) Push(pid PID, r io.Reader) (err error) {
	v := &url.Values{}
	pid.InParameter(v)
	//query url
	u := &url.URL{
		Path:     PUSH,
		RawQuery: v.Encode(),
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	var client http.Client
	req, err := http.NewRequest("POST", c.remote.ResolveReference(u).String(), buf)
	if err != nil {
		return
	}
	req.ContentLength = int64(buf.Len()) // fuck I can't do that, I need to compute the length first
	_, err = client.Do(req)
	return
}

func (c *HttpClient) Search(query string, start int) (result []PID) {
	v := url.Values{}
	v.Set("q", query)
	v.Set("start", strconv.Itoa(start))

	//query url
	u := &url.URL{
		//scheme://[userinfo@]host/path[?query][#fragment]
		Path:     SEARCH,
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(c.remote.ResolveReference(u).String())
	if err != nil {
		return result
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return result
}
