package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tortlewortle/yaigo-svelte/web"
	"github.com/tortlewortle/yaigo/pkg/inertia"
	"github.com/tortlewortle/yaigo/pkg/yaigo"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	inertiaDevMode := os.Getenv("APP_ENV") == "local"

	frontend, err := web.FrontendFS()
	if err != nil {
		log.Fatalf("getting frontend filesystem: %v", err)
	}

	inertiaServer, err := newInertiaServer(frontend, inertiaDevMode)
	if err != nil {
		log.Fatalf("making inertia server: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(inertia.NewMiddleware(inertiaServer))

	r.Handle(
		"/assets/*",
		web.AssetFileServer(frontend),
	)

	r.Get("/", inertia.Handler(func(c *inertia.Ctx, request *http.Request) error {
		return c.Render("Index", inertia.Props{
			"myCoolProp": "balls",
			"mySlowProp": inertia.Defer(func(ctx context.Context) (any, error) {
				time.Sleep(time.Second)
				return "takes a second ok", nil
			}),
		})
	}))

	log.Println("server listening on http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("starting http server: %v", err)
		return
	}
}

func newInertiaServer(frontend fs.FS, useDevServer bool) (*yaigo.Server, error) {
	var opts []yaigo.OptFunc

	if useDevServer {
		opts = append(opts, yaigo.WithViteDevServer("http://localhost:5173", false))
	}

	t := web.RootTemplateFn

	inertiaServer, err := yaigo.NewServer(t, frontend,
		opts...,
	)

	return inertiaServer, err
}
