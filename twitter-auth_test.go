// Twitter Auth - Go

// Copyright 2014 - Dênis Mendes. All rights reserved.
// Author: Dênis Mendes <denisffmendes@gmail.com>
// Use of this source code is governed by a BSD-style

package main

import "testing"

func TestTwitterAuth(t *testing.T) {
  err := twitterAuth()
  if(err != nil){
    t.Error("Error Twitter Auth: %s",err)
  }
}
