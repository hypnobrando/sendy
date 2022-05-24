package sendy

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type (
	// Request is an HTTP builder struct for making an
	// HTTP request.
	Request struct {
		httpClient         *http.Client
		host               string
		hooks              Hooks
		method             string
		path               string
		params             Params
		headers            Headers
		username, password string
		body               io.Reader
		err                error
	}
)

func (request *Request) setErr(err error) *Request {
	request.err = err
	return request
}

// SendIt commits the HTTP request builder and captures the
// response as a Response struct.
func (request *Request) SendIt() *Response {
	var response *Response
	if request.err != nil {
		return response.setErr(request.err)
	}

	params := url.Values{}
	for _, param := range request.params {
		params.Set(param.Key, param.Value)
	}

	var paramString string
	if len(params) > 0 {
		paramString = fmt.Sprintf("?%s", params.Encode())
	}

	url := fmt.Sprintf("%s%s%s", request.host, request.path, paramString)

	httpRequest, err := http.NewRequest(request.method, url, request.body)
	if err != nil {
		return &Response{err: err}
	}

	if request.username != "" || request.password != "" {
		httpRequest.SetBasicAuth(request.username, request.password)
	}

	for _, header := range request.headers {
		httpRequest.Header.Set(header.Key, header.Value)
	}

	for _, hook := range request.hooks {
		hook.Request(httpRequest)
	}

	httpResponse, err := request.httpClient.Do(httpRequest)
	if err != nil {
		return &Response{
			err: err,
		}
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return &Response{err: err}
	}

	return &Response{
		httpResponse: httpResponse,
		body:         body,
		statusCode:   httpResponse.StatusCode,
	}
}

// Method overrides the HTTP request method.
func (request *Request) Method(method string) *Request {
	request.method = method
	return request
}

// Path overrides the HTTP request path.
func (request *Request) Path(path string) *Request {
	request.path = path
	return request
}

// Param appends a URL query param.
func (request *Request) Param(key, value string) *Request {
	request.params = append(request.params, Param{key, value})
	return request
}

// Header appends an HTTP request header.
func (request *Request) Header(key, value string) *Request {
	request.headers = append(request.headers, Header{key, value})
	return request
}

// UserAgent sets the User-Agent request header.
func (request *Request) UserAgent(userAgent string) *Request {
	return request.Header(userAgentHeaderKey, userAgent)
}

// BasicAuth sets the HTTP basic auth username and password.
func (request *Request) BasicAuth(username, password string) *Request {
	request.username = username
	request.password = password
	return request
}

// JSON serializes the input object into the body of the
// request in the form of JSON.
func (request *Request) JSON(object interface{}) *Request {
	if request.err != nil {
		return request
	}

	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return request.setErr(err)
	}

	request.body = bytes.NewReader(jsonBytes)
	return request
}

// XML serializes the input object into the body of the
// request in the form of XML.
func (request *Request) XML(object interface{}) *Request {
	if request.err != nil {
		return request
	}

	xmlBytes, err := xml.Marshal(object)
	if err != nil {
		return request.setErr(err)
	}

	request.body = bytes.NewReader(xmlBytes)
	return request
}

// MultiPartForm serializes the input values into a multi-part form request.
func (request *Request) MultiPartForm(values map[string]io.Reader) *Request {
	if request.err != nil {
		return request
	}

	var b bytes.Buffer
	var err error

	w := multipart.NewWriter(&b)
	defer w.Close()

	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, filepath.Base(x.Name())); err != nil {
				return request.setErr(err)
			}

		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return request.setErr(err)
			}
		}

		if _, err = io.Copy(fw, r); err != nil {
			return request.setErr(err)
		}
	}

	request.body = &b
	return request.Header("Content-Type", w.FormDataContentType())
}

type (
	// Params is a slice of Param structs.
	Params []Param

	// Param contains a single URL query param.
	Param struct {
		Key, Value string
	}
)

type (
	// Headers is a slice of Header structs.
	Headers []Header

	// Header contains a single HTTP header.
	Header struct {
		Key, Value string
	}
)
