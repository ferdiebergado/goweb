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
	r := goexpress.New()
	r.Get("/api/hello", handler.HandleHello)
	req := httptest.NewRequest("GET", "/api/hello", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	const ct = "application/json"
	assert.Equal(t, ct, rr.Header().Get("content-type"))
	apiRes := struct {
		Message string
	}{
		Message: "Hello world!",
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &apiRes); err != nil {
		t.Fatal("failed to decode json", err)
	}

	assert.Equal(t, "Hello world!", apiRes.Message)
}
