package engine

import (
  "bytes"
  "errors"
  h "golang.org/x/net/html"
  "strings"
  "text/template"

  "github.com/robertkrimen/otto"
)

type jses5 struct {
  hasHtml bool
  obj     *otto.Object
  objName string
}

func NewJSES5(src string, objName string) (*jses5, error) {
  vm := otto.New()
  obj, err := vm.Object(src)
  if err != nil {
    return nil, err
  }
  return &jses5{
    false,
    obj,
    objName,
  }, nil
}

func (js *jses5) AddAttr(name string, val interface{}) error {
  if err := js.obj.Set(name, val); err != nil {
    return err
  }
  return nil
}

func (js *jses5) AddMethod(name string, src string) error {
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

func (js *jses5) Build() (*template.Template, error) {
  jsw := newJsWriter(js.objName, "document")

  jsw.affectObj("", "")
  if err := buildInnerObject(js.obj, jsw); err != nil {
    return nil, err
  }

  objToStr := jsw.bf.String()

  if js.hasHtml {
    objToStr = replaceDocQueries(objToStr)
  }

  tmpl := template.Must(template.New("jsObject").Parse(objToStr))

  return tmpl, nil
}

func (js *jses5) GetExt() string {
  return ".min.js"
}

func (js *jses5) IncludeHtml(doc *html) {
  /*
   * buildDoc: function(target){
   *   var _sr = document.querySelector(target).attachShadow({mode:'open'});
   *   // create nodes, add theirs attrs and append children to parents
   *   this.doc = _sr; // To be query in place of the document
   * }
   */

  sRootVar := "_sr" // Shadow root element
  jsw := newJsWriter(sRootVar, "_d")
  jsw.makeFunction("target")
  jsw.affectVar("_d", "document")
  jsw.affectVar(sRootVar, "_d.querySelector(target).attachShadow({mode:'open'})")
  doc.readAndExecute(jsw.buildNode)
  jsw.affectAttr("this", "doc", sRootVar)
  jsw.closeFunction()

  js.AddMethod("buildDoc", jsw.bf.String())

  js.hasHtml = true
}

func (js *jses5) IncludeCss(css string) {
  jsw := newJsWriter("a", "document")

  jsw.makeFunction()
  jsw.affectVar("a", "")
  jsw.createElement("style")
  jsw.affectAttr("a", "innerHTML", css)

  if js.hasHtml {
    jsw.doc = "this.doc"
  }
  jsw.appendChild(jsw.doc, "a")
  jsw.closeFunction()
}


func buildField(field otto.Value, jsw *jsWriter) error {
  var str string
  var err error

  switch {
  case field.IsObject() && !field.IsFunction() && field.Class() != "Array":
    //bf.WriteRune('{')
    jsw.makeObj()
    if err := buildInnerObject(field.Object(), jsw); err != nil {
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
  jsw.bf.WriteString(str)

  return nil
}

func buildInnerObject(obj *otto.Object, jsw *jsWriter) error {
  keys := obj.Keys()

  for i, fieldName := range keys {
    val, err := obj.Get(fieldName)
    if err != nil {
      return err
    }

    // foo: field
    jsw.affectField(fieldName, "")
    if err := buildField(val, jsw); err != nil {
      return err
    }
    if i < len(keys) - 1 {
      jsw.endField()
    }
  }
  jsw.closeObj()

  return nil
}

func formatArray(arr otto.Value) (string, error) {
  val, errExp := arr.Export()
  if errExp != nil {
    return "", errExp
  }
  formatArr, errStr := arr.ToString()
  if errStr != nil {
    return "", errExp
  }
  switch val.(type) {
  case []string:
    var arrBf bytes.Buffer
    split := strings.Split(formatArr, ",")
    for i, val := range split {
      // ["foo","bar"]
      arrBf.WriteRune('"')
      arrBf.WriteString(val)
      arrBf.WriteRune('"')
      if i < len(split) - 1 {
        arrBf.WriteRune(',')
      }
    }
    formatArr = arrBf.String()
  }
  return ("[" + formatArr + "]"), nil
}

func formatVar(pVar otto.Value) (string, error) {
  formatVar, errStr := pVar.ToString()
  if errStr != nil {
    return "", errStr
  }

  switch {
  case pVar.IsString():
    formatVar = "\"" + formatVar + "\""
  case pVar.IsFunction():
    rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
    formatVar = rpcer.Replace(formatVar)
  }

  return formatVar, nil
}

func replaceDocQueries(src string) string {
  rpcer := strings.NewReplacer("document.querySelector", "this.doc.querySelector",
    "document.querySelectorAll", "this.doc.querySelectorAll")
  return rpcer.Replace(src)
}

type jsWriter struct {
  bf    bytes.Buffer
  doc   string
  cVar  string
  vars  []string
}

func newJsWriter(baseVar string, docTarget string) *jsWriter {
  vars := make([]string, 0)
  vars = append(vars, baseVar)
  return &jsWriter{
    bytes.Buffer{},
    docTarget,
    baseVar,
    vars,
  }
}

func (jsw *jsWriter) affectAttr(context string, attrName string, expr string) {
  jsw.bf.WriteString(context)
  jsw.bf.WriteRune('.')
  jsw.bf.WriteString(attrName)
  jsw.bf.WriteString(" = ")
  jsw.bf.WriteString(expr)
  jsw.endExpr()
}

func (jsw *jsWriter) affectField(name string, expr string) {
  jsw.bf.WriteString(name)
  jsw.bf.WriteRune(':')
  jsw.bf.WriteString(expr)
}

func (jsw *jsWriter) affectObj(varName string, obj string) {
  if len(varName) == 0 {
    varName = jsw.cVar
  }
  jsw.bf.WriteString(varName)
  jsw.bf.WriteString("={")
  if len(obj) > 0 {
    jsw.bf.WriteString(obj)
    jsw.bf.WriteRune('}')
    jsw.endExpr()
  }
}

func (jsw *jsWriter) affectVar(varName string, expr string) {
  if len(varName) == 0 {
    varName = jsw.cVar
  }
  jsw.bf.WriteString("var ")
  jsw.bf.WriteString(varName)
  jsw.bf.WriteString(" = ")
  if len(expr) > 0 {
    jsw.bf.WriteString(expr)
    jsw.endExpr()
  }
}

func (jsw *jsWriter) appendChild(to string, toAppend string) {
  if len(toAppend) == 0 {
    toAppend = jsw.cVar
  }
  jsw.bf.WriteString(to)
  jsw.bf.WriteString(".appendChild(")
  jsw.bf.WriteString(toAppend)
  jsw.bf.WriteRune(')')
  jsw.endExpr()
}

func (jsw *jsWriter) buildNode(html *html, index int) {
  switch html.curNode.Type {
  case h.ElementNode:
    if !html.isExcludedNode(html.curNode.Data) {
      jsw.genUniqueVar()
      jsw.affectVar("", "")
      jsw.createElement(html.curNode.Data)
      jsw.setAttributes(html.curNode.Attr)
      jsw.appendChild(jsw.vars[len(jsw.vars) - index - 1], "")
    }
  case h.TextNode:
    jsw.genUniqueVar()
    jsw.affectVar("", "")
    jsw.createTextNode(html.curNode.Data)
    jsw.setAttributes(html.curNode.Attr)
    jsw.appendChild(jsw.vars[len(jsw.vars) - index - 1], "")
  }
}

func (jsw *jsWriter) closeFunction() {
  jsw.bf.WriteRune('}')
}

func (jsw *jsWriter) closeObj() {
  jsw.bf.WriteRune('}')
}

func (jsw *jsWriter) createElement(el string) {
  jsw.bf.WriteString(jsw.doc)
  jsw.bf.WriteString(".createElement(\"")
  jsw.bf.WriteString(el)
  jsw.bf.WriteString("\")")
  jsw.endExpr()
}

func (jsw *jsWriter) createTextNode(text string) {
  jsw.bf.WriteString(jsw.doc)
  jsw.bf.WriteString(".createTextNode(\"")
  jsw.bf.WriteString(text)
  jsw.bf.WriteString("\")")
  jsw.endExpr()
}

func (jsw *jsWriter) endField() {
  jsw.bf.WriteRune(',')
}

func (jsw *jsWriter) endExpr() {
  jsw.bf.WriteRune(';')
}

func (jsw *jsWriter) genUniqueVar() {
  baseNames := [26]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
  tLength := len(jsw.vars)
  bLength := len(baseNames)
  if tLength >= bLength {
    mod := tLength % bLength
    time := tLength / bLength

    jsw.vars = append(jsw.vars, jsw.vars[(time - 1) * bLength] + baseNames[mod])
  } else {
    jsw.vars = append(jsw.vars, baseNames[tLength])
  }
  jsw.cVar = jsw.vars[len(jsw.vars) - 1]
}

func (jsw *jsWriter) setAttributes(attrs []h.Attribute) {
  var attrKey string
  for _, attr := range attrs {
    if len(attr.Namespace) > 0 {
      attrKey = attr.Namespace + ":" + attr.Key
    } else {
      attrKey = attr.Key
    }
    jsw.bf.WriteString(jsw.cVar)
    jsw.bf.WriteString(".setAttribute(\"")
    jsw.bf.WriteString(attrKey)
    jsw.bf.WriteString("\", \"")
    jsw.bf.WriteString(attr.Val)
    jsw.bf.WriteString("\")")
    jsw.endExpr()
  }
}

func (jsw *jsWriter) makeFunction(args ...string) {
  jsw.bf.WriteString("function(")
  for i, arg := range args {
    jsw.bf.WriteString(arg)
    if i < len(args) - 1 {
      jsw.bf.WriteRune(',')
    }
  }
  jsw.bf.WriteString("){")
}

func (jsw *jsWriter) makeObj() {
  jsw.bf.WriteRune('{')
}
