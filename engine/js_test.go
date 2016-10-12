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

var js, _ = engine.NewJS(src, "objForTest")

func TestBuild(t *testing.T) {
  var current, expected bytes.Buffer
  expectedStr := `objForTest={testObj:{field:"yu",fieldObj:{childNum:2.5,childStr:"15 - 15"}},testArrNum:[1,2,3,4],addStr:"hello",addNum:10.2,addArrStr:{0:"str1",1:"str2"},addFn:function(a, b){ console.log('new fn'); }}`

  if err := js.AddAttr("addStr", "hello"); err != nil {
    t.Errorf("AddAttr failed to add a string field, error : %s", err)
  }
  if err := js.AddAttr("addNum", 10.2); err != nil {
    t.Errorf("AddAttr failed to add a number field, error : %s", err)
  }
  if err := js.AddAttr("addArrStr", []string{"str1", "str2"}); err != nil {
    t.Errorf("AddAttr failed to add a string array field, error : %s", err)
  }

  if err := js.AddMethod("addFn", "function(a, b){ console.log('new fn'); }"); err != nil {
    t.Errorf("AddMethod failed to add a function field, error : %s", err)
  }

  tmpl, err := js.Build()
  if err != nil {
    t.Errorf("Build failed, error : %s", err)
  }
  tmplExp := template.Must(template.New("exp").Parse(expectedStr))
  tmpl.Execute(&current, tmpl.Name())
  tmplExp.Execute(&expected, tmplExp.Name())

  t.Log(current.String())
  t.Log(expected.String())

  if current.String() != expected.String() {
    t.Fail()
  }
}
