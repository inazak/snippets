package main

import (
  "golang.org/x/text/encoding/japanese"
  "golang.org/x/text/transform"
  "bytes"
  "io"
  "io/ioutil"
  "os"
  "fmt"
  "strings"
)

var keyword = `そのファイルに含まれているはずの単語`

func main() {

  if len(os.Args) != 2 {
    fmt.Fprintf(os.Stderr, "USAGE: search_word_in_sjis_or_utf8_file FILE")
    os.Exit(1)
  }
  filename := os.Args[1]

  enc, err := DetectEncodingFromKeyword(filename, keyword)
  if err != nil {
    fmt.Fprintf(os.Stderr, err.Error())
    os.Exit(1)
  }

  fmt.Printf("file encoding is %s\n", enc)
}


var encodings = []string{
  "SJIS",
  "UTF8",
}

// determine the encoding based on whether the keyword is included in the file or not.
func DetectEncodingFromKeyword(filename, keyword string) (string, error) {

  for _, enc := range encodings {
    text, err := GetTextFromFile(filename, enc)
    if err != nil {
      return "", fmt.Errorf("error occurred while reading file")
    }
    if index := strings.Index(text, keyword) ; index != -1 {
      return enc, nil
    }
  }

  return "", fmt.Errorf("not contain keyword or unknown encoding file")
}

// reads a file with the specified encoding and returns the text
func GetTextFromFile(filename, encoding string) (string, error) {

  file, err := os.Open(filename)
  if err != nil {
    return "", err
  }
  defer file.Close()

  var reader io.Reader

  switch encoding {
  case "sjis", "SJIS", "ShiftJIS":
    reader = transform.NewReader(file, japanese.ShiftJIS.NewDecoder())

  case "utf8", "UTF8", "UTF-8":
    reader = file

  default:
    return "", fmt.Errorf("unknown encoding")
  }

  done, err := ioutil.ReadAll(reader)
  if err != nil {
    return "", err
  }

  if HasBOM(done) {
    done = StripBOM(done)
  }

  return string(done), nil
}

func HasBOM(in []byte) bool {
	return bytes.HasPrefix(in, []byte{239, 187, 191})
}

func StripBOM(in []byte) []byte {
	return bytes.TrimPrefix(in, []byte{239, 187, 191})
}


