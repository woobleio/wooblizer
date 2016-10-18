package main

import (
  "os"

  "github.com/woobleio/wooblizer/engine"
)

type DocLang int
const (
  HTML DocLang = iotad
)

type ScriptLang int
const (
  JSES5 ScriptLang = iota
)

type Wbzr struct {
  Doc     engine.Doc
  Script  engine.Script
}

func New(sl ScriptLang, scriptSrc string, dl DocLang, docSrc string, name string) (*Wbzr, error) {
  var wbzr Wbzr;

  switch sl {
  case JSES5:
    wbzr.Script, err = engine.NewJSES5(scriptSrc, name)
    if err != nil {
      return nil, err
    }
  default:
    panic("Script not supported")
  }

  switch dl {
  case HTML:
    wbzr.Doc, err = engine.NewHTML(docSrc)
    if err != nil {
      return nil, err
    }
  }

  return &wbzr, nil;
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
