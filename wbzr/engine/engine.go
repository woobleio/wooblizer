package engine

import (
  "text/template"
)

type Engine interface {
  AddMethod(name string, src string)
  Build() (src *template.Template, err error)
  CheckSource(src interface{}) error
  GetExt() string
}
