package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ferdiebergado/gopherkit/http/response"
	"github.com/ferdiebergado/goweb/internal/config"
)

const suffix = ".html"

type templateMap map[string]*template.Template

type Template struct {
	templates templateMap
}

func NewTemplate(cfg config.TemplateConfig) *Template {
	layoutFile := filepath.Join(cfg.Path, cfg.LayoutFile)
	layoutTmpl := template.Must(template.New("layout").Funcs(funcMap()).ParseFiles(layoutFile))
	parsePartials(cfg.Path, cfg.PartialsPath, layoutTmpl)

	return &Template{
		templates: parsePages(cfg.Path, cfg.PagesPath, layoutTmpl),
	}
}

func (t *Template) Render(w http.ResponseWriter, r *http.Request, name string, data any) {
	tmpl, ok := t.templates[name]
	if !ok {
		response.ServerError(w, r, fmt.Errorf("template does not exist: %s", name))
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		response.ServerError(w, r, fmt.Errorf("execute template: %w", err))
		return
	}

	_, err := buf.WriteTo(w)

	if err != nil {
		response.ServerError(w, r, fmt.Errorf("write response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Parse all partial templates into the layout template
func parsePartials(templateDir, partialsDir string, layoutTmpl *template.Template) {
	err := fs.WalkDir(os.DirFS(templateDir), partialsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, suffix) {
			_, parseErr := layoutTmpl.ParseFiles(filepath.Join(templateDir, path))
			if parseErr != nil {
				return fmt.Errorf("parse partials: %w", parseErr)
			}
			slog.Debug("parsed partial", "path", path)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("load partials templates: %w", err))
	}

	slog.Debug("layout", "name", layoutTmpl.Name())
}

// Parse main templates from pagesDir
func parsePages(templateDir, pagesDir string, layoutTmpl *template.Template) templateMap {
	tmplMap := make(templateMap)
	err := fs.WalkDir(os.DirFS(templateDir), pagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, suffix) {
			name := strings.TrimPrefix(path, pagesDir+"/")
			name = strings.TrimSuffix(name, suffix)
			tmplMap[name] = template.Must(template.Must(layoutTmpl.Clone()).ParseFiles(filepath.Join(templateDir, path)))
			slog.Debug("parsed page", "path", path, "name", name)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("load pages templates: %w", err))
	}

	return tmplMap
}

// Retrieve the template func maps
func funcMap() template.FuncMap {
	return template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s) // #nosec G203 -- No user input
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s) // #nosec G203 -- No user input
		},
		"url": func(s string) template.URL {
			return template.URL(s) // #nosec G203 -- No user input
		},
		"js": func(s string) template.JS {
			return template.JS(s) // #nosec G203 -- No user input
		},
		"jsstr": func(s string) template.JSStr {
			return template.JSStr(s) // #nosec G203 -- No user input
		},
		"css": func(s string) template.CSS {
			return template.CSS(s) // #nosec G203 -- No user input
		},
	}
}
