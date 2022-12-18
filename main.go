package main

import (
	"html/template"
	"log"
	"net/http"
)

type templateData struct {
	Path string
  Docs string
  Code string
}

const docRedirect = "https://godoc.formulabun.club/pkg/go.formulabun.club"
const codeRedirect = "https://github.com/formulabun"

const htmlPage = `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="go.formulabun.club{{.Path}} git {{.Code}}{{.Path}}">
<meta http-equiv="refresh" content="0; url={{.Docs}}{{.Path}}">
</head>
<body>
<a href="{{.Docs}}{{.Path}}">Redirecting to documentation.</a>
</body>
</html>`

var templ *template.Template

func handler(w http.ResponseWriter, r *http.Request) {
	data := templateData{
    r.URL.EscapedPath(),
    docRedirect,
    codeRedirect,
  }
	err := templ.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var err error
	templ, err = template.New("redirect").Parse(htmlPage)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8001", nil))
}
