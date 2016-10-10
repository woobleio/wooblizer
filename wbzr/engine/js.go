package engine

import (
  "bytes"
  "strings"
  "text/template"

  "github.com/robertkrimen/otto"
)

type JS struct {
  Obj     *otto.Object
  ObjName string
}

func NewJS(src string, objName string) (*JS, error) {
  vm := otto.New()
  obj, err := vm.Object(src)
  if err != nil {
    return nil, err
  }
  return &JS{obj, objName}, nil
}

func (js *JS) AddMethod(name string, src string) {
}

func (js *JS) Build() (*template.Template, error) {
  var buildBf bytes.Buffer

  // obj = {...}
  buildBf.WriteString(js.ObjName)
  buildBf.WriteString("={")
  if err := parseObject(js.Obj, &buildBf); err != nil {
    return nil, err
  }

  tmpl := template.Must(template.New("jsObject").Parse(buildBf.String()))

  return tmpl, nil
}

func (js *JS) CheckSource(src interface{}) error {
  _, _, err := otto.Run(src)
  return err
}

func (js *JS) GetExt() string {
  return ".min.js"
}

func parseArray(arr otto.Value) (string, error) {
  val, errExp := arr.Export()
  if errExp != nil {
    return "", errExp
  }
  parsedArr, errStr := arr.ToString()
  if errStr != nil {
    return "", errExp
  }
  switch val.(type) {
  case []string:
    var arrBf bytes.Buffer
    split := strings.Split(parsedArr, ",")
    for i, val := range split {
      // ["foo","bar"]
      arrBf.WriteRune('"')
      arrBf.WriteString(val)
      arrBf.WriteRune('"')
      if i < len(split) - 1 {
        arrBf.WriteRune(',')
      }
    }
    parsedArr = arrBf.String()
  }
  return ("[" + parsedArr + "]"), nil
}

func parseField(field otto.Value, bf *bytes.Buffer) error {
  var str string
  var err error

  switch {
  case field.IsObject() && !field.IsFunction() && field.Class() != "Array":
    bf.WriteRune('{')
    if err := parseObject(field.Object(), bf); err != nil {
      return err
    }
  case field.Class() == "Array":
    if str, err = parseArray(field); err != nil {
      return err
    }
  default:
    if str, err = parseVar(field); err != nil {
      return err
    }
  }
  bf.WriteString(str)

  return nil
}

func parseObject(obj *otto.Object, bf *bytes.Buffer) error {
  keys := obj.Keys()

  for i, fieldName := range keys {
    val, err := obj.Get(fieldName)
    if err != nil {
      return err
    }

    // foo: field
    bf.WriteString(fieldName)
    bf.WriteRune(':')
    if err := parseField(val, bf); err != nil {
      return err
    }
    if i < len(keys) - 1 {
      bf.WriteRune(',')
    }
  }
  bf.WriteRune('}')

  return nil
}

func parseVar(pVar otto.Value) (string, error) {
  val, errExp := pVar.Export()
  if errExp != nil {
    return "", errExp
  }
  parsedVar, errStr := pVar.ToString()
  if errStr != nil {
    return "", errStr
  }

  switch val.(type) {
  case string:
    parsedVar = "\"" + parsedVar + "\""
  case map[string]interface{}: // function
    rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
    parsedVar = rpcer.Replace(parsedVar)
  }

  return parsedVar, nil
}
