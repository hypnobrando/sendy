package sendy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingHook struct {
	beforeCount int
	afterCount  int
	methods     []string
	statusCodes []int
	afterErr    error
	headerKey   string
	headerValue string
}

func (h *recordingHook) BeforeRequest(ctx context.Context, req *http.Request) (context.Context, error) {
	h.beforeCount++
	h.methods = append(h.methods, req.Method)
	if h.headerKey != "" {
		req.Header.Set(h.headerKey, h.headerValue)
	}

	return context.WithValue(ctx, hookCtxKey{}, "hooked"), nil
}

func (h *recordingHook) AfterRequest(ctx context.Context, req *http.Request, resp *http.Response, err error) {
	h.afterCount++
	h.afterErr = err
	if resp != nil {
		h.statusCodes = append(h.statusCodes, resp.StatusCode)
	}

	if got, _ := ctx.Value(hookCtxKey{}).(string); got != "hooked" {
		panic("after hook missing before-hook context")
	}
}

type hookCtxKey struct{}

func TestRequestHookSeesRequestAndResponse(t *testing.T) {
	hook := &recordingHook{
		headerKey:   "X-Test-Hook",
		headerValue: "1",
	}

	var sawHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawHeader = r.Header.Get("X-Test-Hook")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	err := NewClient().
		Host(server.URL).
		RequestHook(hook).
		Post().
		Path("/widgets").
		SendIt().
		Error()

	require.NoError(t, err)
	assert.Equal(t, 1, hook.beforeCount)
	assert.Equal(t, 1, hook.afterCount)
	assert.Equal(t, []string{http.MethodPost}, hook.methods)
	assert.Equal(t, []int{http.StatusCreated}, hook.statusCodes)
	assert.Equal(t, "1", sawHeader)
	assert.NoError(t, hook.afterErr)
}

func TestAddRequestHookAppliesToPackageLevelHelpers(t *testing.T) {
	t.Cleanup(resetGlobalRequestHooks)

	hook := &recordingHook{}
	AddRequestHook(hook)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := Get(server.URL).
		Path("/").
		SendIt().
		Error()

	require.NoError(t, err)
	assert.Equal(t, 1, hook.beforeCount)
	assert.Equal(t, 1, hook.afterCount)
}

func TestRequestHookAfterSeesTransportError(t *testing.T) {
	hook := &recordingHook{}

	err := NewClient().
		Host("http://127.0.0.1:1").
		RequestHook(hook).
		Get().
		Path("/").
		SendIt().
		Error()

	require.Error(t, err)
	assert.Equal(t, 1, hook.beforeCount)
	assert.Equal(t, 1, hook.afterCount)
	assert.Error(t, hook.afterErr)
}

func TestRequestHookAfterRunsInReverseOrder(t *testing.T) {
	var order []string

	first := sequenceHook{name: "first", order: &order}
	second := sequenceHook{name: "second", order: &order}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := NewClient().
		Host(server.URL).
		RequestHook(first).
		RequestHook(second).
		Get().
		SendIt().
		Error()

	require.NoError(t, err)
	assert.Equal(t, []string{"first-before", "second-before", "second-after", "first-after"}, order)
}

type sequenceHook struct {
	name  string
	order *[]string
}

func (h sequenceHook) BeforeRequest(ctx context.Context, req *http.Request) (context.Context, error) {
	*h.order = append(*h.order, h.name+"-before")
	return ctx, nil
}

func (h sequenceHook) AfterRequest(ctx context.Context, req *http.Request, resp *http.Response, err error) {
	*h.order = append(*h.order, h.name+"-after")
}
