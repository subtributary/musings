package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/subtributary/musings/internal/localization"
)

type Responder struct {
	views *ViewFactory
}

func NewResponder(liveTemplates bool, locales []localization.Locale, contentRoot *os.Root, staticRoot *os.Root) (Responder, error) {
	views, err := NewViewFactory(liveTemplates, locales, contentRoot, staticRoot)
	if err != nil {
		return Responder{}, fmt.Errorf("new view factory: %w", err)
	}
	return Responder{views: views}, nil
}

func (resp *Responder) File(w http.ResponseWriter, r *http.Request, root *os.Root, name string) {
	file, err := root.Open(name)
	if err != nil {
		resp.NotFound(w, r)
		return
	}
	defer (func() { _ = file.Close() })()

	info, err := file.Stat()
	if err != nil {
		resp.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, name, info.ModTime(), file)
}

func (resp *Responder) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	resp.View(w, r, "404")
}

func (resp *Responder) ServerError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type ViewOption func(*View) error

func WithData(data any) ViewOption {
	return func(v *View) error {
		return v.SetData(data)
	}
}

func (resp *Responder) View(w http.ResponseWriter, r *http.Request, name string, opts ...ViewOption) {
	view, err := resp.views.Create(r, name)
	if err != nil {
		resp.ServerError(w, fmt.Errorf("create view: %w", err))
		return
	}

	for _, opt := range opts {
		if err = opt(view); err != nil {
			resp.ServerError(w, fmt.Errorf("configure view: %w", err))
			return
		}
	}

	if err = view.Serve(w); err != nil {
		resp.ServerError(w, fmt.Errorf("serve view: %w", err))
	}
}
