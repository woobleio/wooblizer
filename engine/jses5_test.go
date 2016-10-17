package engine

import (
  "bytes"
  "testing"
  "text/template"

  "github.com/woobleio/wooblizer/engine"
)

var src =
`
({
  testObj: {
    field: "yu",
    fieldObj: {
      childNum: 2.5,
      childStr: "15 - 15"
    }
  },
  testArrNum: [1, 2, 3, 4]
})
`

var js, _ = engine.NewJSES5(src, "objForTest")

func TestBuild(t *testing.T) {
  var current, expected bytes.Buffer
  expectedStr := `objForTest={testObj:{field:"yu",fieldObj:{childNum:2.5,childStr:"15 - 15"}},testArrNum:[1,2,3,4],addStr:"hello",addNum:10.2,addArrStr:{0:"str1",1:"str2"},addFn:function(a, b){ console.log('new fn'); document.querySelector('#elid'); document.querySelectorAll('div'); document.querySelectorAll('.elclass'); }}`

  if err := js.AddAttr("addStr", "hello"); err != nil {
    t.Errorf("AddAttr failed to add a string field, error : %s", err)
  }
  if err := js.AddAttr("addNum", 10.2); err != nil {
    t.Errorf("AddAttr failed to add a number field, error : %s", err)
  }
  if err := js.AddAttr("addArrStr", []string{"str1", "str2"}); err != nil {
    t.Errorf("AddAttr failed to add a string array field, error : %s", err)
  }

  if err := js.AddMethod("addFn", "function(a, b){ console.log('new fn'); document.querySelector('#elid'); document.querySelectorAll('div'); document.querySelectorAll('.elclass'); }"); err != nil {
    t.Errorf("AddMethod failed to add a function field, error : %s", err)
  }

  tmpl, err := js.Build()
  if err != nil {
    t.Errorf("Build failed, error : %s", err)
  }
  tmplExp := template.Must(template.New("exp").Parse(expectedStr))
  tmpl.Execute(&current, tmpl.Name())
  tmplExp.Execute(&expected, tmplExp.Name())

  if current.String() != expected.String() {
    t.Fail()
  }

  current.Reset()
  expected.Reset()

  expectedStr = `objForTest={testObj:{field:"yu",fieldObj:{childNum:2.5,childStr:"15 - 15"}},testArrNum:[1,2,3,4],addStr:"hello",addNum:10.2,addArrStr:{0:"str1",1:"str2"},addFn:function(a, b){ console.log('new fn'); this.doc.querySelector('#elid'); this.doc.querySelectorAll('div'); this.doc.querySelectorAll('.elclass'); },buildDoc:function(target){var _d = document;var _sr = _d.querySelector(target).attachShadow({mode:'open'});var b = _d.createElement("div");b.setAttribute("class", "classel");_sr.appendChild(b);var c = _d.createElement("p");c.setAttribute("id", "paragraph");b.appendChild(c);var d = _d.createTextNode("this is a text");c.appendChild(d);var e = _d.createElement("div");e.setAttribute("data", "a data");d.appendChild(e);var f = _d.createElement("span");f.setAttribute("class", "first-class second-class");f.setAttribute("id", "spanid");e.appendChild(f);this.doc = _sr;}}`

  docHtml, err := engine.NewHTML("<div class='classel'><p id='paragraph'>this is a text</p><div data='a data'></div></div><span class='first-class second-class' id='spanid'></span>")
  if err != nil {
    t.Errorf("NewHTML failed, error : %s", err)
  }
  docHtml.AddExcludedNodes("body", "html", "head")
  js.IncludeHtml(docHtml)
  tmpl, err = js.Build()
  if err != nil {
    t.Errorf("Build failed, error : %s", err)
  }
  tmplExp = template.Must(template.New("exp").Parse(expectedStr))
  tmpl.Execute(&current, tmpl.Name())
  tmplExp.Execute(&expected, tmplExp.Name())

  if current.String() != expected.String() {
    t.Fail()
  }
}
