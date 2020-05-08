package sendy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Response contains the response after making an HTTP
// request and can be used for parsing / inspecting
// the response.
type Response struct {
	httpResponse *http.Response
	err          error
}

// StatusCode returns the HTTP status code of the response.
// If there was an error making the request and no response was
// read then the status code returned is -1.
func (response *Response) StatusCode() int {
	if response.httpResponse != nil {
		return response.httpResponse.StatusCode
	}

	return -1
}

func (response *Response) setErr(err error) *Response {
	response.err = err
	return response
}

// Close completes reading the HTTP response and closes the connection.
// Close must be called if none of the other response reading
// functions are called, otherwise there will be a leak in
// HTTP connections.
func (response *Response) Close() *Response {
	if response.err != nil {
		return response
	}

	_, err := ioutil.ReadAll(response.httpResponse.Body)
	if err != nil {
		return response.setErr(err)
	}

	err = response.httpResponse.Body.Close()
	return response.setErr(err)
}

// JSON parses the response body as JSON and deserializes it
// into the input object.
func (response *Response) JSON(object interface{}) *Response {
	if response.err != nil {
		return response
	}

	err := json.NewDecoder(response.httpResponse.Body).Decode(object)
	if err != nil {
		return response.setErr(err)
	}

	return response.Close()
}

// Error returns an Error struct which can be accepted as
// an error interface.  This Error contains an error's that
// might have occurred during the build process, during
// the lifetime of the request, or even during the parsing of
// the response.  Non-200 status codes are also returned
// as an Error.
func (response *Response) Error() *Error {
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
