package script

import (
  "github.com/robertkrimen/otto"
)

type JS struct {
}

func (js *JS) TranspileObject(src string) (*WbObject, error) {
  vm := otto.New()
  //return vm.Object(src) Format obj
  return nil
}
