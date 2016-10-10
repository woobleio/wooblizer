package wbzr

import (
  "os"

  "github.com/woobleio/wooblizer/wbzr/engine"
)

// Implemented engines
const (
  JS int = iota
)

type Wbzr struct {
  Engine  engine.Engine
}

func New(eng int, src string, name string) *Wbzr {
  var wbzr Wbzr;

  switch(eng) {
  case JS:
    wbzr.Engine, _ = engine.NewJS(src, name)
  default:
    panic("Engine not supported")
  }

  return &wbzr;
}

func (wb *Wbzr) BuildFile(path string, fileName string) error {
  tmplSrc, err := wb.Engine.Build()
  if err != nil {
    return err
  }

  f, err := os.Create(path + "/" + fileName + wb.Engine.GetExt())
  if err != nil {
    return err
  }

  defer f.Close()

  if err := tmplSrc.Execute(f, tmplSrc.Name()); err != nil {
    return err
  }

  return nil
}
