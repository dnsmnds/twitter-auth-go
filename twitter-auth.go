// Twitter Auth - Go

// Copyright 2014 - Dênis Mendes. All rights reserved.
// Author: Dênis Mendes <denisffmendes@gmail.com>
// Use of this source code is governed by a BSD-style

package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

type TwitterAuthError struct {
  Errors []json.RawMessage `json:"errors"`
}

type ErrorInfo struct {
  Code int8
  Label string
  Message string
}

func main() {
  url := "https://api.twitter.com/oauth2/token";

  resp, err := http.Get(url)
  if err != nil {
    fmt.Printf("Error: %s", err)
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  twitterAuthErrors := &TwitterAuthError{}
  err = json.Unmarshal(body, &twitterAuthErrors)

  if err != nil {
    fmt.Printf("Error: %s \n\n",err)
  }

  for _, t := range twitterAuthErrors.Errors {
    errorInfo := &ErrorInfo{}
    err := json.Unmarshal(t, &errorInfo)
    if err != nil {
      fmt.Println("%s",err)
    }
    fmt.Printf("%v %s %s \n\n",errorInfo.Code, errorInfo.Label, errorInfo.Message)
  }

}

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
