package sendy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type (
	// Hooks is a slice of Hook interfaces.
	Hooks []Hook

	// Hook is used as a hook before a Client makes
	// a request.
	Hook interface {
		Request(*http.Request)
	}
)

// Hook takes in a Hook interface that gets
// called right before an HTTP request is made.
func (c *Client) Hook(hook Hook) *Client {
	c.hooks = append(c.hooks, hook)
	return c
}

// Hook takes in a Hook interface that gets
// called right before an HTTP request is made.
func (request *Request) Hook(hook Hook) *Request {
	request.hooks = append(request.hooks, hook)
	return request
}

// DumpRequests prints out the entire contents of every
// request to stdout.
func (c *Client) DumpRequests() *Client {
	return c.Hook(&dumpRequestHook{})
}

// Dump prints out the entire contents of the request
// when the request is made.
func (r *Request) Dump() *Request {
	return r.Hook(&dumpRequestHook{})
}

type dumpRequestHook struct{}

// Request implements the Hook interface for the
// dumpRequestHook.  This simply dumps the entire contents
// of the request to stdout.
func (hook *dumpRequestHook) Request(request *http.Request) {
	rawRequest, _ := httputil.DumpRequest(request, true)
	fmt.Println("\n" + string(rawRequest) + "\n")
}
