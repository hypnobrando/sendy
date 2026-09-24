package sendy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"
)

type (
	// Hooks is a slice of Hook interfaces.
	Hooks []Hook

	// Hook is used as a hook before a Client makes
	// a request.
	Hook interface {
		Request(*http.Request)
	}

	// RequestHooks is a slice of RequestHook interfaces.
	RequestHooks []RequestHook

	// RequestHook observes an HTTP request before it is sent
	// and after it completes. BeforeRequest may return a
	// derived context (for example an OpenTelemetry span)
	// that is attached to the request.
	RequestHook interface {
		BeforeRequest(ctx context.Context, req *http.Request) (context.Context, error)
		AfterRequest(ctx context.Context, req *http.Request, resp *http.Response, err error)
	}
)

var (
	globalRequestHooksMu sync.RWMutex
	globalRequestHooks   RequestHooks
)

// AddRequestHook registers a RequestHook that is applied to
// every request, including package-level helpers such as Get
// and Post. Hooks are read when the request is sent so they
// apply to clients created earlier.
func AddRequestHook(hook RequestHook) {
	globalRequestHooksMu.Lock()
	defer globalRequestHooksMu.Unlock()

	globalRequestHooks = append(globalRequestHooks, hook)
}

func copyGlobalRequestHooks() RequestHooks {
	globalRequestHooksMu.RLock()
	defer globalRequestHooksMu.RUnlock()

	if len(globalRequestHooks) == 0 {
		return nil
	}

	hooks := make(RequestHooks, len(globalRequestHooks))
	copy(hooks, globalRequestHooks)

	return hooks
}

func resetGlobalRequestHooks() {
	globalRequestHooksMu.Lock()
	defer globalRequestHooksMu.Unlock()

	globalRequestHooks = nil
}

// Hook takes in a Hook interface that gets
// called right before an HTTP request is made.
func (c *Client) Hook(hook Hook) *Client {
	c.hooks = append(c.hooks, hook)
	return c
}

// RequestHook takes in a RequestHook that is called
// before and after the HTTP request is sent.
func (c *Client) RequestHook(hook RequestHook) *Client {
	c.requestHooks = append(c.requestHooks, hook)
	return c
}

// Hook takes in a Hook interface that gets
// called right before an HTTP request is made.
func (request *Request) Hook(hook Hook) *Request {
	request.hooks = append(request.hooks, hook)
	return request
}

// RequestHook takes in a RequestHook that is called
// before and after the HTTP request is sent.
func (request *Request) RequestHook(hook RequestHook) *Request {
	request.requestHooks = append(request.requestHooks, hook)
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
