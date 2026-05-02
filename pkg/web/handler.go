package web

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nzwice/surl/pkg/web/templates"
	"github.com/templui/templui/utils"
)

//go:embed templates/*
var templatesFS embed.FS

var (
	parsedTemplates = template.Must(template.ParseFS(templatesFS, "templates/*"))
)

var (
	pages = map[string]templ.Component{
		"index": templates.Index(),
	}
)

func Static(staticDir string) http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))
}

func RegisterTemplUIScripts(mux *http.ServeMux, debug bool) {
	utils.SetupScriptRoutes(mux, debug)
}

func Page(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p, ok := pages[name]; ok {
			templ.Handler(templates.Page(p)).ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	})
}
