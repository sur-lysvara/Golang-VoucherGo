package main

import (
	"html/template"
	"net/http"
)

func publicPage(w http.ResponseWriter, title string, body template.HTML) {
	t := template.Must(template.New("public-minimal-layout").Parse(publicPageLayout))

	if err := t.Execute(w, map[string]any{
		"Title": title,
		"Body":  body,
	}); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
