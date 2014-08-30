// Twitter Auth - Go

// Copyright 2014 - Dênis Mendes. All rights reserved.
// Author: Dênis Mendes <denisffmendes@gmail.com>
// Use of this source code is governed by a BSD-style

package main

import (
  "bytes"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "encoding/base64"
  "strconv"
  "compress/gzip"
  "io"
  "net/url"
)

//-----------------------------------------------------------------------------
// Constants
//-----------------------------------------------------------------------------

const (
  URL_OAUTH2 = "https://api.twitter.com/oauth2/token"
  CONSUMER_KEY = "MY_CONSUMER_KEY"
  CONSUMER_SECRET = "MY_CONSUMER_SECRET"
)

//-----------------------------------------------------------------------------
// Structs
//-----------------------------------------------------------------------------

type TwitterAuthError struct {
  //-------------------------------------------
  // RawMessage is a raw encoded JSON object.
  // Ref.: goo.gl/JicizA
  //-------------------------------------------
  Errors []json.RawMessage `json:"errors"`
}

type TwitterAccessToken struct {
  TokenType string `json:"token_type"`
  AccessToken string `json:"access_token"`
}

type ErrorInfo struct {
  Code int
  Label string
  Message string
}

type OAuth2Credentials struct {
  Authorization string
  ContentType string
  ContentLength int
  AcceptEncoding string
  GrantType string
}

//-----------------------------------------------------------------------------

func main() {
  twitterAuth()
}

func twitterAuth() (err error) {

  //----------------------------------------------------------------------------
  // EncodeToString returns the base64 encoding of src.
  // Ref.: goo.gl/D8p05W
  //----------------------------------------------------------------------------
  authorizationEncoded := base64.StdEncoding.EncodeToString([]byte(CONSUMER_KEY + ":"+ CONSUMER_SECRET))
  //----------------------------------------------------------------------------

  //----------------------------------------------------------------------------
  // Configuration
  //----------------------------------------------------------------------------
  oAuth2Credentials := OAuth2Credentials{
    "Basic " + authorizationEncoded,
    "application/x-www-form-urlencoded;charset=UTF-8",
    29,
    "gzip",
    "client_credentials",
  }
  //----------------------------------------------------------------------------

  //----------------------------------------------------------------------------
  // Http Client
  // Ref.: goo.gl/M67Fd0
  //----------------------------------------------------------------------------
  client := &http.Client{}

  urlValues := url.Values{}
  urlValues.Add("grant_type", "client_credentials");

  req, err := http.NewRequest("POST",URL_OAUTH2, bytes.NewBufferString(urlValues.Encode()))
  if err != nil {
    fmt.Printf("Error http.NewRequest: %s", err)
  }
  req.Header.Add("Authorization", oAuth2Credentials.Authorization)
  req.Header.Add("Content-Type", oAuth2Credentials.ContentType)
  req.Header.Add("Content-Length", strconv.Itoa(oAuth2Credentials.ContentLength))
  req.Header.Add("Accept-Encoding", oAuth2Credentials.AcceptEncoding)
  req.Header.Add("User-Agent", "Twitter Auth Go")
  resp, err := client.Do(req);
  if err != nil {
    fmt.Printf("Error client.Do: %s", err)
  }
  defer resp.Body.Close()
  //----------------------------------------------------------------------------

  //----------------------------------------------------------------------------
  // Response body
  //----------------------------------------------------------------------------
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("Error ioutil.ReadAll: %s \n\n",err)
  }
  //----------------------------------------------------------------------------

  if parseErrorAuth(body) {
    return
  }

  parseAccessToken(body)

  return
}

func generateReader(body []byte) (io.Reader){
  //----------------------------------------------------------------------------
  // Decompress the content into an io.Reader
  // Ref.: goo.gl/vNfZKR
  //----------------------------------------------------------------------------
  var reader io.Reader
  var buffer *bytes.Buffer
  var err error

  buffer = bytes.NewBuffer(body)
  reader, err = gzip.NewReader(buffer)
  if err != nil {
    fmt.Printf("Error gzip.NewReader: %s \n\n",err)
    return nil
  }
  //----------------------------------------------------------------------------

  return reader
}

func parseAccessToken(body []byte) (error){
  //----------------------------------------------------------------------------
  // Use the stream interface to decode json from the io.Reader
  // Ref.: goo.gl/vNfZKR
  //----------------------------------------------------------------------------
  twitterAccessToken := new (TwitterAccessToken);
  err := json.NewDecoder(generateReader(body)).Decode(twitterAccessToken);
  if err != nil && err != io.EOF {
    fmt.Println("Error json.NewDecoder(reader).Decode: %s",err)
    return err
  }
  //----------------------------------------------------------------------------

  fmt.Printf("TokenType: %s \nAccessToken: %s \n",twitterAccessToken.TokenType,twitterAccessToken.AccessToken)

  return nil
}

func parseErrorAuth(body []byte) (error bool) {
  //----------------------------------------------------------------------------
  // Use the stream interface to decode json from the io.Reader
  // Ref.: goo.gl/vNfZKR
  //----------------------------------------------------------------------------
  var twitterAuthErrors = &TwitterAuthError{};
  dec := json.NewDecoder(generateReader(body))
  err := dec.Decode(twitterAuthErrors)
  if err != nil && err != io.EOF {
    fmt.Println("Error dec.Decode: %s",err)
    return true
  }
  //----------------------------------------------------------------------------

  for _, t := range twitterAuthErrors.Errors {
    errorInfo := &ErrorInfo{}
    err := json.Unmarshal(t, &errorInfo)
    if err != nil {
      fmt.Println("Error json.Unmarshal: %s",err)
      return true
    }
    fmt.Printf("Code: %v \nLabel: %s \nMessage: %s \n\n",errorInfo.Code, errorInfo.Label, errorInfo.Message)
  }

  if len(twitterAuthErrors.Errors) > 0 {
    return true
  }

  return false
}

//-----------------------------------------------------------------------------
// Template - Tests
//-----------------------------------------------------------------------------

func test(){
  /*
  {
    "errors": [
        {
            "code": 99,
            "label": "authenticity_token_error",
            "message": "Unable to verify your credentials"
        }
    ]
  }*/

  var jsonLocal = []byte(`
    {"errors": [ {"code": 99, "label": "authenticity_token_error", "message": "Unable to verify your credentials"} ] }
  `)

  twitterAuthErrors := &TwitterAuthError{}
  err := json.Unmarshal(jsonLocal, &twitterAuthErrors)

  if err != nil {
    fmt.Printf("Error: %s \n\n",err)
  }

  for _, t := range twitterAuthErrors.Errors {
    dst := &ErrorInfo{}
    err := json.Unmarshal(t, &dst)
    if err != nil {
      fmt.Println("%s",err)
    }
    fmt.Printf("%s \n\n",dst.Message)
  }
}
