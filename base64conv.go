package main

import (
  "fmt"
  "encoding/base64"
  "io/ioutil"
  "os"
)

func isExist(filename string) bool {
  _, err := os.Stat(filename)
  return err == nil
}

func main() {
  if len(os.Args) != 3 {
    fmt.Fprintf(os.Stderr, "USAGE: base64conv enc|dec file.txt")
    os.Exit(1)
  }

  filename := os.Args[2]
  tempfile := filename + ".tmp"

  if ! isExist(filename) {
    fmt.Fprintf(os.Stderr, "USAGE: base64conv enc|dec file.txt")
    os.Exit(1)
  }

  switch os.Args[1] {
  case "enc":
    os.Remove(tempfile)
    os.Rename(filename, tempfile)
    dat, _ := ioutil.ReadFile(tempfile)
    enc := base64.StdEncoding.EncodeToString(dat)
    ioutil.WriteFile(filename, []byte(enc), os.ModePerm)
    os.Remove(tempfile)
  case "dec":
    os.Remove(tempfile)
    os.Rename(filename, tempfile)
    dat, _ := ioutil.ReadFile(tempfile)
    dec, _ := base64.StdEncoding.DecodeString(string(dat))
    ioutil.WriteFile(filename, dec, os.ModePerm)
    os.Remove(tempfile)
  default:
    fmt.Fprintf(os.Stderr, "USAGE: base64conv enc|dec file.txt")
    os.Exit(1)
  }
}

