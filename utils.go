package main

import (
  "math/rand"
)

var AlphaLU = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(length int, randAlphabet []rune) string {
  stub := make([]rune, length)
  for i := range stub {
    stub[i] = randAlphabet[rand.Intn(len(randAlphabet))]
  }
  return string(stub)
}
