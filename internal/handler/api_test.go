package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleHello(t *testing.T) {
	const url = "/api/hello"
	const msg = "Hello world!"
	const ct = "application/json"

	r := goexpress.New()
	r.Get(url, handler.HandleHello)

	req := httptest.NewRequest("GET", url, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.Equal(t, ct, rr.Header().Get("Content-Type"))

	apiRes := handler.APIResponse{
		Message: msg,
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &apiRes); err != nil {
		t.Fatal("failed to decode json", err)
	}

	assert.Equal(t, msg, apiRes.Message)
}
