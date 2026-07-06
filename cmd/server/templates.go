package main

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var templateFiles embed.FS

type Templates struct {
	Err404 *template.Template
	Index  *template.Template
	Post   *template.Template
}

func LoadTemplates() (Templates, error) {
	err404, err := template.ParseFS(templateFiles, "templates/*.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load 404 template parts: %w", err)
	}
	_, err = err404.ParseFS(templateFiles, "templates/pages/404.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load 404 template: %w", err)
	}
	err404 = err404.Lookup("layout.gohtml")

	index, err := template.ParseFS(templateFiles, "templates/*.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load index template parts: %w", err)
	}
	_, err = index.ParseFS(templateFiles, "templates/pages/index.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load index template: %w", err)
	}
	index = index.Lookup("layout.gohtml")

	post, err := template.ParseFS(templateFiles, "templates/*.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load post template parts: %w", err)
	}
	_, err = post.ParseFS(templateFiles, "templates/pages/post.gohtml")
	if err != nil {
		return Templates{}, fmt.Errorf("load post template: %w", err)
	}
	post = post.Lookup("layout.gohtml")

	return Templates{
		Err404: err404,
		Index:  index,
		Post:   post,
	}, nil

}
