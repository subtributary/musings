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

type Response struct {
	r *http.Request
	w http.ResponseWriter
}

func newResponse(w http.ResponseWriter, r *http.Request) Response {
	return Response{r: r, w: w}
}

func (r Response) File(root *os.Root, name string) {
	http.ServeFileFS(r.w, r.r, root.FS(), name)
}

func (r Response) NotFound() {
	http.Error(r.w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (r Response) ServerError(err error) {
	log.Printf("server error: %v", err)
	http.Error(r.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type JSONResponse struct {
	Response
}

func NewJSONResponse(w http.ResponseWriter, r *http.Request) JSONResponse {
	return JSONResponse{
		Response: newResponse(w, r),
	}
}

func (r JSONResponse) Okay(data any) {
	r.w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(r.w).Encode(data)
	if err != nil {
		r.ServerError(err)
	}
}

type ViewResponse struct {
	Response
}

func NewViewResponse(w http.ResponseWriter, r *http.Request) ViewResponse {
	return ViewResponse{
		Response: newResponse(w, r),
	}
}

func (r ViewResponse) Okay(tmpl *template.Template, data any) {
	r.view(tmpl, data)
}

func (r ViewResponse) NotFound(tmpl *template.Template, data any) {
	r.w.WriteHeader(http.StatusNotFound)
	r.view(tmpl, data)
}

func (r ViewResponse) view(tmpl *template.Template, viewModel any) {
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
