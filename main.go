// This HTTP utility servers a HTML comment containing
// a meta tag named "go-import", and adheres to
// https://go.dev/ref/mod#serving-from-proxy
// This application should be run by piping
// a json configuration document into its
// standard input. Like `goimport < config.json`
// See `go doc goimport Config` for the configuration
// file structure.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

// Configuration for the HTML document containing a meta tag
// redirecting `go get` to the source code of the module
// See: https://go.dev/ref/mod#vcs-find
type Config struct {

	// Root url of all modules, which is used in the go.get command.
	// HTTP GET requests to this url are handled by this application
	// This is the template to printf and is given one
	// string argument containing the request path
	RootPath string `json:"moduleRoot"`
	// Root url of the vcs holding all modules
	// This is the template to printf and is given one
	// string argument containing the request path
	CodeRoot string `json:"vcsRoot"`
	// vcs system, like: git bzr fossil svn ...
	// empty for http GET
	Vcs string `json:"vcs"`
	// Root url of the human facing documentation
	// This is the template to printf and is given one
	// string argument containing the request path
	// Browser requests to RootPath will be redirected
	// to this URL
	DocRoot string `json:"docRoot"`

	// Port to host the http server on
	Port int

	Path string `json:"-"`
}

const htmlPage = `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="{{ printf .RootPath .Path }} {{.Vcs}} {{ printf .CodeRoot .Path }}">
<meta http-equiv="refresh" content="0; url={{ printf .DocRoot .Path }}">
</head>
<body>
<a href="{{ printf .DocRoot .Path }}">Redirecting to documentation.</a>
</body>
</html>`

var templ *template.Template
var config Config

func handler(w http.ResponseWriter, r *http.Request) {
	var c struct{
		Config
		Path string
	}
	c.Config = config
	c.Path = r.URL.EscapedPath()
	if c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}
	err := templ.Execute(w, c)
	if err != nil {
		log.Println(err)
	}
}

func parseConfig() error {
	stdInstat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if (stdInstat.Mode() & fs.ModeCharDevice) != 0 {
		return errors.New("stdin is not a character device")
	}
	err = json.NewDecoder(os.Stdin).Decode(&config)
	if err != nil {
		return fmt.Errorf("Config file is invalid: %w", err)
	}
	return nil
}

func printUsage() {
	fmt.Println(`
Run this utility by piping in a json config file.
The structure of the file is defined by go.openfl.eu/goimport.Config
	`)
}

func main() {
	if err := parseConfig(); err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}
	templ = template.Must(template.New("redirect").Parse(htmlPage))
	http.HandleFunc("/", handler)
	hostAddr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("hosting on %s\n", hostAddr)
	log.Fatal(http.ListenAndServe(hostAddr, nil))
}
