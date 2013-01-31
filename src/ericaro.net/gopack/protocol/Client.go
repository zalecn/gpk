package protocol

import (
	"io"
	"net/url"
)

// Client constructor is used to associate a client to it's url scheme
type ClientConstructor func(name string, u url.URL, token *Token) (c Client, err error)

// We keep track of all client factory in a map url scheme -> ClientConstructor 
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

//NewClient allocates a client (i.e a remote) with the three elements shared by all remotes :
// name, an url, and a Token.
func NewClient(name string, u url.URL, token *Token) (Client, error) {
	//fmt.Printf("new remote %s %v. scheme factory = %s\n", name, u.String(), RemoteRepositoryFactory[u.Scheme])
	return ClientFactory[u.Scheme](name, u, token)
}

//Client is any kind of client that can talk to a remote repository. In the commands, it is called a Remote
type Client interface {
	//Fetch retrieve the package in the form of a io.ReadCloser. 
	//r reads into a tar.gz stream, containing all the packages files, including the .gpk
	Fetch(pid PID) (r io.ReadCloser, err error)
	//Push will send what's in the reader to the remote.
	//r is expected to be a tar.gz stream containing all the packages files, including the .gpk
	Push(pid PID, r io.Reader) (err error)
	
	//PushExecutables send the content of ./bin directory. If the build is cross platform, every binary is expected to be in the right directory
	PushExecutables(pid PID, r io.Reader) (err error)
	
	//Search run a search query on the remote. The exact meaning of "query" is free ( remote can implement operators if they want).
	// start is the offset where the start sending the results.
	// result is a slice of PID returned. The number of returned result is free (usually limited to 10)
	Search(query string, start int) (result []PID)
	//Name the remote's name: the way it is referenced to from the command line. Must be unique
	Name() string
	//Path the remote's URL: any valid URL. Usually clients are bound to an URL scheme. The client can do whatever he wants with it
	Path() url.URL
	//Token the remote's Token : a Token is an authentication Token provided by the remote server to avoid exchanging passwords.
	// usually the token can be discarded on the server. Some permissions are granted to a given Token.
	// There is only one Token per remote
	Token() *Token
}

// A base implementation of a Client, does implement partially the Client interface.
// therefore it can be used as a delegation field in a real client.
type BaseClient struct {
	url url.URL
	name  string
	token *Token
}

//NewBaseClient fill a BaseClient with all it need
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