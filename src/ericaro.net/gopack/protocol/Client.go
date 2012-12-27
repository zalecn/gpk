package protocol

import (
	"io"
	"net/url"
)

type ClientConstructor func(name string, u url.URL, token *Token) (c Client, err error)

var ClientFactory map[string]ClientConstructor // factory
func RegisterClient(urlprotocol string, xtor ClientConstructor) {
	if ClientFactory == nil {
		ClientFactory = make(map[string]ClientConstructor)
	}
	if _, ok := ClientFactory[urlprotocol]; ok {
		panic("double remote repository definition for " + urlprotocol + "\n")
	}
	ClientFactory[urlprotocol] = xtor
}

func NewClient(name string, u url.URL, token *Token) (Client, error) {
	//fmt.Printf("new remote %s %v. scheme factory = %s\n", name, u.String(), RemoteRepositoryFactory[u.Scheme])
	return ClientFactory[u.Scheme](name, u, token)
}

//Client is any kind of client that can talk to a remote repository
type Client interface {
	Fetch(pid PID) (r io.ReadCloser, err error)
	Push(pid PID, r io.Reader) (err error)
	Search(query string, start int) (result []PID)
	Name() string
	Path() url.URL
	Token() *Token
}

type BaseClient struct {
	url url.URL
	name  string
	token *Token
}

func NewBaseClient(name string, u url.URL, token *Token) (r *BaseClient) {
	r = &BaseClient{
		name: name,
		token: token,
		url: u,
	}
	return

}

func (r BaseClient) Token() *Token { return r.token }
func (r BaseClient) Name() string { return r.name }
func (r BaseClient) Path() url.URL { return r.url }