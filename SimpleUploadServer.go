package main

import (
  "html/template"
  "io"
  "net/http"
  "os"
  "os/exec"
  "path/filepath"
  "fmt"
)

var (
  fileSavePath = ""
  Action = ""
)

func saveHandler(w http.ResponseWriter, r *http.Request) {
  filetype := r.FormValue("filetype")

  data, _, err := r.FormFile("file")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer data.Close()

  file, err := os.Create(filepath.Join(fileSavePath,filetype))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Fprintf(w, "%s", "error: file create")
    return
  }
  defer file.Close()

  _, errio := io.Copy(file, data)
  if errio != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Fprintf(w, "%s", "error: file copy")
    return
  }

  if Action != "" {
    exec.Command(Action, filetype).Run()
  }

  http.Redirect(w, r, "/upload", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
  var templ = template.Must(template.ParseFiles("upload.html"))
  templ.Execute(w, "upload.html")
}

func main() {
  if len(os.Args) < 3 {
    fmt.Println(`Usage: SimpleUploadServer PORT FILE_SAVE_PATH [ACTION]`)
    fmt.Println(`Description:`)
    fmt.Println(`  http server use ./upload.html`)
    fmt.Println(`  filename is saved from 'filetype' value`)
    fmt.Println(`  ACTION is execute commandline after saved`)
    fmt.Println(`  --- upload HTML snippet ---`)
    fmt.Println(`  <form method="post" action="/save" enctype="multipart/form-data">`)
    fmt.Println(`  <input type="hidden" name="filetype" value="nameToSave">`)
    fmt.Println(`  <input type="file" name="file">`)
    fmt.Println(`  <input type="submit" name="submit" value="UPLOAD">`)
    fmt.Println(`  </form>`)
    os.Exit(1)
  }

  port := os.Args[1]
  fileSavePath = os.Args[2]
  if len(os.Args) == 4 {
    Action = os.Args[3]
  }

  http.HandleFunc("/upload", uploadHandler)
  http.HandleFunc("/save", saveHandler)

  http.ListenAndServe(":"+port, nil)
}

