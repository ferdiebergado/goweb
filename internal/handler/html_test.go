package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/pkg/template"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleDashboard(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()

	mockCfg := config.TemplateConfig{
		Path:         "../../web/templates",
		LayoutFile:   "layout.html",
		PartialsPath: "partials",
		PagesPath:    "pages",
	}
	tmpl := template.NewTemplate(mockCfg)
	h := handler.NewBaseHTMLHandler(tmpl)

	r := goexpress.New()
	r.Get("/dashboard", h.HandleDashboard)
	r.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code should match")
	assert.Contains(t, rr.Body.String(), "Dashboard", "Body should contain the same text")
}
