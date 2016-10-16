package main

import (
  "os"

  "github.com/woobleio/wooblizer/engine"
)

// Implemented engines
const (
  JS int = iota
)

type Wbzr struct {
  Script  engine.Script
}

func New(langTarget int, src string, name string) *Wbzr {
  var wbzr Wbzr;

  switch(langTarget) {
  case JS:
    wbzr.Script, _ = engine.NewJSES5(src, name)
  default:
    panic("Script not supported")
  }

  return &wbzr;
}

func (wb *Wbzr) BuildFile(path string, fileName string) error {
  tmplSrc, err := wb.Script.Build()
  if err != nil {
    return err
  }

  f, err := os.Create(path + "/" + fileName + wb.Script.GetExt())
  if err != nil {
    return err
  }

  defer f.Close()

  if err := tmplSrc.Execute(f, tmplSrc.Name()); err != nil {
    return err
  }

  return nil
}
