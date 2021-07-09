package main

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"
)

func main() {

  pwd, _ := os.Getwd()
  filter := filepath.Join(pwd, "*")

  filelist, err := GetFilepathList(filter)
  if err != nil {
    fmt.Printf("can't create a list of files")
    os.Exit(1)
  }

  for _, filename := range filelist {
    fmt.Printf("%s\n", filename)
  }
}


func GetFilepathList(filter string) (list []string, err error) {

  entries, err := filepath.Glob(filter)
  if err != nil {
    return nil, err
  }

  for _, entry := range entries {
    isdir, err := isDir(entry)
    if err != nil {
      return nil, err
    }
    if ! isdir {
      list = append(list, entry)
    }
  }

  return list, nil
}

func isDir(p string) (bool, error) {
  f, err := os.Stat(p)
  if err != nil {
    return false, err
  }
  return f.Mode().IsDir(), nil
}

func splitFilepath(fpath string) (dir, name, ext string) {

  dir   = filepath.Dir(fpath)
  base := filepath.Base(fpath)

  dotindex := strings.LastIndex(base, ".")

  // file has no extention
  if dotindex == -1 {
    name = base
    ext  = ""
    return
  }

  name = base[0:dotindex]
  ext  = base[dotindex:len(base)]
  return
}

