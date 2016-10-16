package engine

import (
  "bytes"
  "errors"
  "golang.org/x/net/html"
  "strings"
  "text/template"

  "github.com/robertkrimen/otto"
)

type JSES5 struct {
  ObjName string

  hasHtml bool
  obj     *otto.Object
}

func NewJSES5(src string, objName string) (*JSES5, error) {
  vm := otto.New()
  obj, err := vm.Object(src)
  if err != nil {
    return nil, err
  }
  return &JSES5{objName, false, obj}, nil
}

func (js *JSES5) AddAttr(name string, val interface{}) error {
  if err := js.obj.Set(name, val); err != nil {
    return err
  }
  return nil
}

func (js *JSES5) AddMethod(name string, src string) error {
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

  if !fn.IsFunction() {
    return errors.New("The source should be a JS function(...params){}")
  }
  if err := js.obj.Set(name, fn); err != nil {
    return err
  }

  return nil
}

func (js *JSES5) Build() (*template.Template, error) {
  var buildBf bytes.Buffer

  // obj = {...}
  buildBf.WriteString(js.ObjName)
  buildBf.WriteString("={")
  if err := writeInnerObject(js.obj, &buildBf); err != nil {
    return nil, err
  }

  objToStr := buildBf.String()

  if js.hasHtml {
    objToStr = replaceDocQueries(objToStr)
  }

  tmpl := template.Must(template.New("jsObject").Parse(objToStr))

  return tmpl, nil
}

func (js *JSES5) GetExt() string {
  return ".min.js"
}

func (js *JSES5) IncludeHtml(doc string) error {
  var bf bytes.Buffer

  r := strings.NewReader(doc)
  html, err := html.Parse(r)
  if err != nil {
    return err
  }

  /*
   * buildDoc: function(target){
   *   var _sr = document.querySelector(target).attachShadow({mode:'open'});
   *   // create nodes, add theirs attrs and append children to parents
   *   this.doc = _sr; // To be query in place of the document
   * }
   */
  sRootVar := "_sr" // Shadow root element
  bf.WriteString("function(target){")
  bf.WriteString("var _d = document;")
  bf.WriteString("var ")
  bf.WriteString(sRootVar)
  bf.WriteString(" = _d.querySelector(target).attachShadow({mode:'open'});") // TODO for ES6 -> https://developers.google.com/web/fundamentals/getting-started/primers/shadowdom#slots
  writeNodesBuilder(html, &bf, &[]string{}, sRootVar)
  bf.WriteString("this.doc = ")
  bf.WriteString(sRootVar)
  bf.WriteRune('}')

  js.AddMethod("buildDoc", bf.String())

  js.hasHtml = true

  return nil
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
  case pVar.IsFunction():
    rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
    parsedVar = rpcer.Replace(parsedVar)
  }

  return parsedVar, nil
}

func getUniqueVar(t *[]string) string {
  baseNames := [26]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
  tLength := len(*t)
  bLength := len(baseNames)
  if tLength >= bLength {
    mod := tLength % bLength
    time := tLength / bLength

    *t = append(*t, (*t)[(time - 1) * bLength] + baseNames[mod])
  } else {
    *t = append(*t, baseNames[tLength])
  }
  return (*t)[len(*t) - 1]
}

func isExcludedNode(node string) bool {
  exclNodes := []string{"html", "body", "head", "script"}
  for _, n := range exclNodes {
    if n == node {
      return true
    }
  }
  return false
}

func jsAppendChild(toVar string, nodeVar string) string {
  return toVar + ".appendChild(" + nodeVar + ");"
}

func jsCreateElement(varName string, el string) string {
  return "var " + varName + " = _d.createElement('" + el + "');"
}

func jsCreateTextNode(varName string, text string) string {
  return "var " + varName + " = _d.createTextNode(\"" + text + "\");"
}

func jsSetAttribute(varName string, attr string, value string, namespace string) string {
  if len(namespace) > 0 {
    attr = namespace + ":" + attr
  }
  return varName + ".setAttribute('" + attr + "', '" + value + "');"
}

func replaceDocQueries(src string) string {
  rpcer := strings.NewReplacer("document.querySelector", "this.doc.querySelector",
    "document.querySelectorAll", "this.doc.querySelectorAll")
  return rpcer.Replace(src)
}

func writeField(field otto.Value, bf *bytes.Buffer) error {
  var str string
  var err error

  switch {
  case field.IsObject() && !field.IsFunction() && field.Class() != "Array":
    bf.WriteRune('{')
    if err := writeInnerObject(field.Object(), bf); err != nil {
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

func writeInnerObject(obj *otto.Object, bf *bytes.Buffer) error {
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

func writeNodesBuilder(n *html.Node, bf *bytes.Buffer, vars *[]string, parentVar string) {
  currentVar, createNode := "", ""

  switch n.Type {
  case html.ElementNode:
    if !isExcludedNode(n.Data) {
      currentVar = getUniqueVar(vars)
      createNode = jsCreateElement(currentVar, n.Data)
    }
  case html.TextNode:
    currentVar = getUniqueVar(vars)
    createNode = jsCreateTextNode(currentVar, n.Data)
  }

  if len(currentVar) == 0 {
    currentVar = parentVar
  }
  if len(createNode) > 0 {
    bf.WriteString(createNode)
    for _, attr := range n.Attr {
      bf.WriteString(jsSetAttribute(currentVar, attr.Key, attr.Val, attr.Namespace))
    }
    bf.WriteString(jsAppendChild(parentVar, currentVar))
  }

  for c := n.FirstChild; c !=nil; c = c.NextSibling {
    writeNodesBuilder(c, bf, vars, currentVar)
  }
}
