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

func (js *JS) AddAttr(name string, val interface{}) error {
  if err := js.Obj.Set(name, val); err != nil {
    return err
  }
  return nil
}

func (js *JS) AddMethod(name string, src string) error {
  vm := otto.New()

  // TODO this is a workaround to build a fn with Otto
  tmpObj, err := vm.Object("({tmp:" + src + "})")
  if err != nil {
    return err
  }
  fn, err := tmpObj.Get("tmp")
  if err != nil {
    return err
  }

  if err := js.Obj.Set(name, fn); err != nil {
    return err
  }

  return nil
}

func (js *JS) Build() (*template.Template, error) {
  var buildBf bytes.Buffer

  // obj = {...}
  buildBf.WriteString(js.ObjName)
  buildBf.WriteString("={")
  if err := writeObject(js.Obj, &buildBf); err != nil {
    return nil, err
  }

  tmpl := template.Must(template.New("jsObject").Parse(buildBf.String()))

  return tmpl, nil
}

func (js *JS) CheckSource(src interface{}) error {
  // TODO care for attacks
  _, _, err := otto.Run(src)
  return err
}

func (js *JS) GetExt() string {
  return ".min.js"
}

func formatArray(arr otto.Value) (string, error) {
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


func formatVar(pVar otto.Value) (string, error) {
  parsedVar, errStr := pVar.ToString()
  if errStr != nil {
    return "", errStr
  }

  switch {
  case pVar.IsString():
    parsedVar = "\"" + parsedVar + "\""
  case pVar.IsFunction(): // function
    rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
    parsedVar = rpcer.Replace(parsedVar)
  }

  return parsedVar, nil
}

func writeField(field otto.Value, bf *bytes.Buffer) error {
  var str string
  var err error

  switch {
  case field.IsObject() && !field.IsFunction() && field.Class() != "Array":
    bf.WriteRune('{')
    if err := writeObject(field.Object(), bf); err != nil {
      return err
    }
  case field.Class() == "Array":
    if str, err = formatArray(field); err != nil {
      return err
    }
  default:
    if str, err = formatVar(field); err != nil {
      return err
    }
  }
  bf.WriteString(str)

  return nil
}

func writeObject(obj *otto.Object, bf *bytes.Buffer) error {
  keys := obj.Keys()

  for i, fieldName := range keys {
    val, err := obj.Get(fieldName)
    if err != nil {
      return err
    }

    // foo: field
    bf.WriteString(fieldName)
    bf.WriteRune(':')
    if err := writeField(val, bf); err != nil {
      return err
    }
    if i < len(keys) - 1 {
      bf.WriteRune(',')
    }
  }
  bf.WriteRune('}')

  return nil
}
