package handler_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleHello(t *testing.T) {
	r := goexpress.New()
	r.Get("/hello", handler.HandleHello)
	req := httptest.NewRequest("GET", "/hello", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, rr.Body.String(), "Hello world!")
}
