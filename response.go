package sendy

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Response contains the response after making an HTTP
// request and can be used for parsing / inspecting
// the response.
type Response struct {
	body       []byte
	statusCode int
	err        error
}

// StatusCode returns the HTTP status code of the response.
// If there was an error making the request and no response was
// read then the status code returned is -1.
func (response *Response) StatusCode() int {
	if response.statusCode != 0 {
		return response.statusCode
	}

	return -1
}

func (response *Response) setErr(err error) *Response {
	response.err = err
	return response
}

// JSON parses the response body as JSON and deserializes it
// into the input object.
func (response *Response) JSON(object interface{}) *Response {
	if response.err != nil {
		return response
	}

	err := json.NewDecoder(bytes.NewReader(response.body)).Decode(object)
	if err != nil {
		return response.setErr(err)
	}

	return response
}

// XML parses the response body as XML and deserializes it
// into the input object.
func (response *Response) XML(object interface{}) *Response {
	if response.err != nil {
		return response
	}

	err := xml.NewDecoder(bytes.NewReader(response.body)).Decode(object)
	if err != nil {
		return response.setErr(err)
	}

	return response
}

// Raw returns the raw body and any errors associated with the request.
func (response *Response) Raw() ([]byte, error) {
	return response.body, response.err
}

// Error returns an Error struct.  This Error contains an error's that
// might have occurred during the build process, during
// the lifetime of the request, or even during the parsing of
// the response.  Non-2XX status codes are also returned
// as an Error.
func (response *Response) Error() error {
	if response.err != nil {
		return &Error{
			err: response.err,
		}
	}

	statusCode := response.StatusCode()
	if statusCode >= 300 {
		return &Error{
			statusCode: statusCode,
		}
	}

	return nil
}

// Headers returns the response headers.
func (response *Response) Headers() http.Header {
	return response.Headers()
}
