package engine

import (
  "text/template"
)

type Script interface {
  AddAttr(name string, val interface{}) error
  AddMethod(name string, src string) error
  Build() (src *template.Template, err error)
  GetExt() string
  IncludeHtml(html string) error
}
