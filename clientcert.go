package main

// クライアント証明書を使ったHTTPSサイトへの接続

import (
  "fmt"
  "crypto/tls"
  "crypto/x509"
  "net/http"
  "io/ioutil"
)

func main() {

  caPem := []byte(`-----BEGIN CERTIFICATE-----
#####PASTE#####
-----END CERTIFICATE-----`)

  certPem := []byte(`-----BEGIN CERTIFICATE-----
#####PASTE#####
-----END CERTIFICATE-----`)

  keyPem := []byte(`-----BEGIN PRIVATE KEY-----
#####PASTE#####
-----END PRIVATE KEY-----`)


  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caPem)

  cert, err := tls.X509KeyPair(certPem, keyPem)
  if err != nil {
    fmt.Printf("%v", err)
    return
  }

  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{ cert },
    RootCAs: caCertPool,
    Renegotiation: tls.RenegotiateOnceAsClient,
    InsecureSkipVerify: true,
  }

  tlsConfig.BuildNameToCertificate()

  transport := &http.Transport{ TLSClientConfig: tlsConfig }
  client := &http.Client{ Transport: transport }

  // --------------------------

  resp, err := client.Get("https://server:8080/xxx/xxxx/index.html")
  if err != nil {
    fmt.Printf("%v", err)
    return
  }
  defer resp.Body.Close()

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("%v", err)
    return
  }

  fmt.Printf("%v", string(data))
}


