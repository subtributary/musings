package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type ResponseOption func(r *Response)

func WithData(data any) ResponseOption {
	return func(r *Response) {
		r.data = data
	}
}

func WithView(view *template.Template) ResponseOption {
	return func(r *Response) {
		r.view = view
	}
}

type Response struct {
	r    *http.Request
	w    http.ResponseWriter
	view *template.Template
	data any
}

func NewResponse(w http.ResponseWriter, r *http.Request) Response {
	return Response{r: r, w: w}
}

func (r *Response) File(root *os.Root, name string) {
	http.ServeFileFS(r.w, r.r, root.FS(), name)
}

func (r *Response) NotFound(options ...ResponseOption) {
	r.applyOptions(options)

	if r.view != nil {
		r.w.WriteHeader(http.StatusNotFound)
		r.writeView()
		return
	}

	http.Error(r.w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (r *Response) Okay(options ...ResponseOption) {
	r.applyOptions(options)

	if r.view != nil {
		r.writeView()
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(r.w).Encode(r.data); err != nil {
		r.ServerError(err)
	}
}

func (r *Response) ServerError(err error) {
	log.Printf("server error: %v", err)
	http.Error(r.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r *Response) writeView() {
	// Write to a buffer so that errors do not leave it partially written.
	var buf bytes.Buffer
	err := r.view.Execute(&buf, r.data)
	if err != nil {
		err = fmt.Errorf("execute template: %w", err)
		r.ServerError(err)
		return
	}

	r.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(r.w)
	if err != nil {
		err = fmt.Errorf("write response: %w", err)
		r.ServerError(err)
		return
	}
}

func (r *Response) applyOptions(options []ResponseOption) {
	for _, opt := range options {
		opt(r)
	}
}
