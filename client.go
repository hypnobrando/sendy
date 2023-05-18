package sendy

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

type (
	// Client is a struct used for making HTTP requests to a
	// a single host.
	Client struct {
		host               string
		username, password string
		hooks              Hooks
		headers            []Header
		httpClient         *http.Client
	}
)

// NewClient instantiates a client for making HTTP requests
// to a single host.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func defaultClient(host string) *Client {
	return NewClient().
		Host(host).
		WithRetriesAndTimeout(3, 12*time.Second)
}

// WithRetriesAndTimeout sets the HTTP transport on the client.
func (c *Client) WithRetriesAndTimeout(maxRetries int, timeout time.Duration) *Client {
	return c.
		Transport(defaultTransport(maxRetries, timeout)).
		Timeout(time.Duration(maxRetries) * timeout)
}

func defaultTransport(maxRetries int, timeout time.Duration) *rehttp.Transport {
	return rehttp.NewTransport(
		nil,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(maxRetries),
			rehttp.RetryAny(
				rehttp.RetryTemporaryErr(),
				rehttp.RetryIsErr(func(err error) bool {
					if err == nil {
						return false
					}

					return strings.Contains(err.Error(), "net/http: request canceled (Client.Timeout exceeded while awaiting headers)")
				}),
			),
		),
		rehttp.ExpJitterDelay(time.Second, timeout),
	)
}

// Host overrides the host of the Client.
func (c *Client) Host(host string) *Client {
	c.host = host
	return c
}

// Transport sets the HTTP transport on the client.  This
// incldues proxying, retry policy, etc.
func (c *Client) Transport(transport http.RoundTripper) *Client {
	c.httpClient.Transport = transport
	return c
}

// Timeout sets the client request timeout.
func (c *Client) Timeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// Header appends an HTTP request header.
func (c *Client) Header(key, value string) *Client {
	c.headers = append(c.headers, Header{key, value})
	return c
}

const userAgentHeaderKey = "User-Agent"

// UserAgent sets the User-Agent HTTP Header.
func (c *Client) UserAgent(userAgent string) *Client {
	return c.Header(userAgentHeaderKey, userAgent)
}

// BasicAuth sets the HTTP basic auth username and password.
func (c *Client) BasicAuth(username, password string) *Client {
	c.username = username
	c.password = password
	return c
}

func (c *Client) request() *Request {
	return &Request{
		httpClient: c.httpClient,
		host:       c.host,
		headers:    c.headers,
		username:   c.username,
		password:   c.password,
		hooks:      c.hooks,
	}
}

// Get instantiates an HTTP GET request builder.
func Get(host string) *Request {
	return defaultClient(host).Get()
}

// Get instantiates an HTTP GET request builder.
func (c *Client) Get() *Request {
	return c.request().Method(http.MethodGet)
}

// Post instantiates an HTTP Post request builder.
func Post(host string) *Request {
	return defaultClient(host).Post()
}

// Post instantiates an HTTP Post request builder.
func (c *Client) Post() *Request {
	return c.request().Method(http.MethodPost)
}

// Patch instantiates an HTTP Patch request builder.
func Patch(host string) *Request {
	return defaultClient(host).Patch()
}

// Patch instantiates an HTTP Patch request builder.
func (c *Client) Patch() *Request {
	return c.request().Method(http.MethodPatch)
}

// Put instantiates an HTTP Put request builder.
func Put(host string) *Request {
	return defaultClient(host).Put()
}

// Put instantiates an HTTP Put request builder.
func (c *Client) Put() *Request {
	return c.request().Method(http.MethodPut)
}

// Delete instantiates an HTTP Delete request builder.
func Delete(host string) *Request {
	return defaultClient(host).Delete()
}

// Delete instantiates an HTTP Delete request builder.
func (c *Client) Delete() *Request {
	return c.request().Method(http.MethodDelete)
}

// Head instantiates an HTTP Head request builder.
func Head(host string) *Request {
	return defaultClient(host).Head()
}

// Head instantiates an HTTP Head request builder.
func (c *Client) Head() *Request {
	return c.request().Method(http.MethodHead)
}

// Connect instantiates an HTTP Connect request builder.
func Connect(host string) *Request {
	return defaultClient(host).Connect()
}

// Connect instantiates an HTTP Connect request builder.
func (c *Client) Connect() *Request {
	return c.request().Method(http.MethodConnect)
}

// Trace instantiates an HTTP Trace request builder.
func Trace(host string) *Request {
	return defaultClient(host).Trace()
}

// Trace instantiates an HTTP Trace request builder.
func (c *Client) Trace() *Request {
	return c.request().Method(http.MethodTrace)
}

// Options instantiates an HTTP Options request builder.
func Options(host string) *Request {
	return defaultClient(host).Options()
}

// Options instantiates an HTTP Options request builder.
func (c *Client) Options() *Request {
	return c.request().Method(http.MethodOptions)
}
