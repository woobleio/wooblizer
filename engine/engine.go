package engine

import (
  "text/template"
)

type Engine interface {
  AddAttr(name string, val interface{}) error
  AddMethod(name string, src string) error
  Build() (src *template.Template, err error)
  CheckSource(src interface{}) error
  GetExt() string
}
