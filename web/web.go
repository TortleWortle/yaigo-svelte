package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed index.gohtml
var rootTemplateHTML string

//go:embed dist/*
var dist embed.FS

func FrontendFS() (fs.FS, error) {
	frontend, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}
	return frontend, nil
}

func RootTemplateFn(t *template.Template) (*template.Template, error) {
	return t.Parse(rootTemplateHTML)
}

func AssetFileServer(frontend fs.FS) http.Handler {
	fs := http.FileServer(http.FS(frontend))
	return fs
}
