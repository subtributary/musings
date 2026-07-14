package web

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Response struct {
	r *http.Request
	w http.ResponseWriter
}

func NewResponse(w http.ResponseWriter, r *http.Request) Response {
	return Response{r: r, w: w}
}

func (r Response) File(root *os.Root, name string) {
	http.ServeFileFS(r.w, r.r, root.FS(), name)
}

func (r Response) NotFound(notFoundTemplate *template.Template, viewModel any) {
	r.w.WriteHeader(http.StatusNotFound)
	r.View(notFoundTemplate, viewModel)
}

func (r Response) ServerError(err error) {
	log.Printf("server error: %v", err)
	http.Error(r.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r Response) View(tmpl *template.Template, viewModel any) {
	// Write to a buffer so that errors do not leave it partially written.
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, viewModel)
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
