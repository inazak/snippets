package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "regexp"
  "strings"
  "time"
)

var interval time.Duration = 45

//NOTE: golang 1.8 or later must use url.PathEscape
func pathEscape(s string) string {
  r := url.QueryEscape(s)
  r = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(r, "$1%20")
  return r
}

func makeGoogleSearchURL(keyword string) string {
  base := "https://www.google.co.jp/search?"
  v := url.Values{}
  v.Set("q", keyword)
  return base + v.Encode()
}

func getHTML(url string) string {
  client  := &http.Client{}
  req, _  := http.NewRequest("GET", url, nil)
  req.Header.Set("User-Agent", "Mozilla/5.0")
  r, _ := client.Do(req)
  defer r.Body.Close()
  b, _ := ioutil.ReadAll(r.Body)
  return string(b)
}

func getImage(url, saveas string) {
  client  := &http.Client{}
  req, _  := http.NewRequest("GET", url, nil)
  req.Header.Set("User-Agent", "Mozilla/5.0")
  r, _ := client.Do(req)
  defer r.Body.Close()
  file, _ := os.Create(saveas)
  defer file.Close()
  io.Copy(file, r.Body)
  return
}


var re_amazonurl, _ = regexp.Compile(`href="/url\?q=(https://www.amazon.co.jp/[^/]+/dp/([^?&]+))`)
var re_amazonimg, _ = regexp.Compile(`image-stretch-vertical frontImage[^;]+;(https://images-na.ssl-images-amazon.com/images/[^&]+)`)

func extractAmazonURL(html string) (url, asin string, err error) {
  match := re_amazonurl.FindStringSubmatch(html)
  if match == nil {
    return "", "", fmt.Errorf("amazon url not found")
  }
  return match[1], match[2], nil
}

func extractAmazonImageURL(html string) (url string, err error) {
  match := re_amazonimg.FindStringSubmatch(html)
  if match == nil {
    return "", fmt.Errorf("image url not found")
  }
  return match[1], nil
}


func main() {

  // title.txt is book-title-list by UTF-8 and LF
  txt, err := ioutil.ReadFile("title.txt")
  if err != nil {
    panic(err)
  }

  // titles by line
  title := strings.Split(string(txt), "\n")

  for _, t := range title {
    if t == "" { continue }

    p := pathEscape(t)
    r := "" // infomation result

    // google search
    searchurl := makeGoogleSearchURL(t) + "+1" //search first volume
    searchhtm := getHTML(searchurl)

    // parse html
    amazonurl, asin, err := extractAmazonURL(searchhtm)

    if err != nil {
      // save infomation
      r = fmt.Sprintf("title=%v\nescaped=%v\ngoogle=%v\n", t, p, searchurl)
      ioutil.WriteFile(p + ".info.txt", []byte(r), os.ModePerm)
      continue
    }

    // get amazon page
    amazonhtm := getHTML(amazonurl)
    ioutil.WriteFile(p + ".amazon.html", []byte(amazonhtm), os.ModePerm)

    // parse html
    imageurl, err := extractAmazonImageURL(amazonhtm)

    if err != nil {
      // save infomation
      r = fmt.Sprintf("title=%v\nescaped=%v\ngoogle=%v\namazon=%v\nasin=%v\n", t, p, searchurl, amazonurl, asin)
      ioutil.WriteFile(p + ".info.txt", []byte(r), os.ModePerm)
      continue
    }

    // save infomation
    r = fmt.Sprintf("title=%v\nescaped=%v\ngoogle=%v\namazon=%v\nasin=%v\nimageurl=%v\n", t, p, searchurl, amazonurl, asin, imageurl)
    ioutil.WriteFile(p + ".info.txt", []byte(r), os.ModePerm)

    // save cover image
    getImage(imageurl, p + ".jpg")

    time.Sleep(interval * time.Second)
  }
}


