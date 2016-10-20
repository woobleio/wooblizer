package main

import (
  _ "os"

  "github.com/woobleio/wooblizer/engine"
  "github.com/woobleio/wooblizer/engine/script"
)

type ScriptLang int
const (
  JSES5 ScriptLang = iota
)

type wbzr struct {
  lang     ScriptLang
  scripts  map[string]engine.Script
}

func New(sl ScriptLang) *wbzr {
  switch sl {
  case JSES5:
  default:
    panic("Script not supported")
  }

  return &wbzr{sl, make(map[string]engine.Script)};
}

func (wb *wbzr) BuildFile(path string, fileName string) error {
  /*tmplSrc, err := wb.script.Build()
  if err != nil {
    return err
  }

  f, err := os.Create(path + "/" + fileName + wb.script.GetExt())
  if err != nil {
    return err
  }
  defer f.Close()

  if err := tmplSrc.Execute(f, tmplSrc.Name()); err != nil {
    return err
  }*/

  return nil
}

func (wb *wbzr) Inject(obj string, name string) (*engine.Script, error) {
  var sc engine.Script
  switch wb.lang {
  case JSES5:
    sc, err := script.NewJSES5(obj, name)
    if err != nil {
      return nil, err
    }
    wb.scripts[name] = sc
  }
  return &sc, nil
}
