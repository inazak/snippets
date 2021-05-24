package main

import (
  "html/template"
  "net/http"
  "fmt"
)

const port = "22222"

const html = `
<html><head>SimpleInputOutputHttpServer</head><body>
<form method="POST" action="/">
<textarea name="input"  cols="80" rows="12">{{.Input}}</textarea><br />
<input type="submit" value="submit" /><br />
<textarea name="output" cols="80" rows="12">{{.Output}}</textarea><br />
</form>
</body></html>
`

var templ, _ = template.New("").Parse(html)

type Page struct {
  Input  string
  Output string
}

func handler(w http.ResponseWriter, r *http.Request) {

  input  := ""
  output := ""

  v := r.FormValue("input")
  if v == "" {
    input  = "enter something"
    output = ""
  } else {
    // do something
    input  = ""
    output = "your input is " + v
  }

  templ.Execute(w, Page{ input, output } )
}

func main() {
  fmt.Printf("the server is starting.\n")
  http.HandleFunc("/", handler)
  http.ListenAndServe(":" + port, nil)
}

