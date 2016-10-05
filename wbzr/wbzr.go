package wbzr

import (
  "github.com/woobleio/wooblizer/wbzr/script"
)

type Wbzr struct {
  script  script.Script
  src     string
}

func New(scriptLanguage string, src string) *Wbzr {
  var wbzr Builder;

  switch(scriptLanguage) {
  case "js":
    wbzr.script = &script.JS{}
  }

  wbzr.src = src

  return &wbzr;
}

func (wb *Wbzr) BuildFile() (bool, error) {
  wb.script.TranspileObject(wb.src)

  return true, nil
}
