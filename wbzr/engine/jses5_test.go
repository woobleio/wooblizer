package engine_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/woobleio/wooblizer/wbzr/engine"
)

func TestBuild(t *testing.T) {
	var current, expected bytes.Buffer

	src :=
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

	js, err := engine.NewJSES5(src, "objForTest")
	if err != nil {
		t.Errorf("Can't create a new jses5, error : %s", err)
	}

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

	/*
	 * objForTest={
	 *   testObj:{field:"yu",fieldObj:{childNum:2.5,childStr:"15 - 15"}},
	 *   testArrNum:[1,2,3,4],
	 *   addStr:"hello",
	 *   addNum:10.2,
	 *   addArrStr:{0:"str1",1:"str2"},
	 *   addFn:function(a, b){
	 *     console.log('new fn');
	 *     document.querySelector('#elid');
	 *     document.querySelectorAll('div');
	 *     document.querySelectorAll('.elclass');
	 *   }
	 * }
	 */
	expectedStr := `objForTest={testObj:{field:"yu",fieldObj:{childNum:2.5,childStr:"15 - 15"}},testArrNum:[1,2,3,4],addStr:"hello",addNum:10.2,addArrStr:{0:"str1",1:"str2"},addFn:function(a, b){ console.log('new fn'); document.querySelector('#elid'); document.querySelectorAll('div'); document.querySelectorAll('.elclass'); }}`

	tmpl, err := js.Build()
	if err != nil {
		t.Errorf("Build failed, error : %s", err)
	}
	tmplExp := template.Must(template.New("exp").Parse(expectedStr))
	tmpl.Execute(&current, js)
	tmplExp.Execute(&expected, tmplExp.Name())

	if current.String() != expected.String() {
		t.Logf("Current is  : %s", current.String())
		t.Logf("Expected is : %s", expected.String())
		t.Fail()
	}
}

func TestIncludeHtml(t *testing.T) {
	var current, expected bytes.Buffer
	js, err := engine.NewJSES5("({})", "objForTest")

	/*
	 * objForTest={
	 *   _builDoc:function(target){
	 *     var _d = document;
	 *     var _sr = _d.querySelector(target).attachShadow({mode:'open'});
	 *     var b = _d.createElement("div");
	 *     b.setAttribute("class", "classel");
	 *     _sr.appendChild(b);
	 *     var c = _d.createElement("p");
	 *     c.setAttribute("id", "paragraph");
	 *     b.appendChild(c);
	 *     var d = _d.createTextNode("this is a text");
	 *     c.appendChild(d);
	 *     var e = _d.createElement("div");
	 *     e.setAttribute("data", "a data");
	 *     b.appendChild(e);
	 *     var f = _d.createElement("span");
	 *     f.setAttribute("class", "first-class second-class");
	 *     f.setAttribute("id", "spanid");
	 *     _sr.appendChild(f);
	 *     this._doc = _sr;
	 *   }
	 * }
	 */
	expectedStr := `objForTest={_buildDoc:function(target){var _d = document;var _sr = _d.querySelector(target).attachShadow({mode:'open'});var b = _d.createElement("div");b.setAttribute("class", "classel");_sr.appendChild(b);var c = _d.createElement("p");c.setAttribute("id", "paragraph");b.appendChild(c);var d = _d.createTextNode("this is a text");c.appendChild(d);var e = _d.createElement("div");e.setAttribute("data", "a data");b.appendChild(e);var f = _d.createElement("span");f.setAttribute("class", "first-class second-class");f.setAttribute("id", "spanid");_sr.appendChild(f);this._doc = function() { return _sr};}}`

	if err := js.IncludeHtml("<div class='classel'><p id='paragraph'>this is a text</p><div data='a data'></div></div><span class='first-class second-class' id='spanid'></span>"); err != nil {
		t.Errorf("IncludeHtml failed, error : %s", err)
	}

	tmpl, err := js.Build()
	if err != nil {
		t.Errorf("Build failed, error : %s", err)
	}
	tmplExp := template.Must(template.New("exp").Parse(expectedStr))

	tmpl.Execute(&current, js)
	tmplExp.Execute(&expected, tmplExp.Name())

	if current.String() != expected.String() {
		t.Logf("Current is  : %s", current.String())
		t.Logf("Expected is : %s", expected.String())
		t.Fail()
	}
}

func TestIncludeCss(t *testing.T) {
	var current, expected bytes.Buffer
	js, err := engine.NewJSES5("({})", "objForTest")

	css :=
		`
  p {
    color: red;
  }
  #id {
    font-size: 2em;
  }
  `

	/*
	 * objForTest={
	 *   _buildStyle:function(){
	 *     var a = document.createElement("style");
	 *     a.innerHTML = "  p {    color: red;  }  #id {    font-size: 2em;  }  ";
	 *     document.appendChild(a);
	 *   }
	 * }
	 */
	expectedStr := `objForTest={_buildStyle:function(){var a = document.createElement("style");a.innerHTML = "  p {    color: red;  }  #id {    font-size: 2em;  }  ";this._doc().appendChild(a);}}`

	if err := js.IncludeCss(css); err != nil {
		t.Errorf("IncludeCss failed, error : %s", err)
	}

	tmpl, err := js.Build()
	if err != nil {
		t.Errorf("Build failed, error : %s", err)
	}
	tmplExp := template.Must(template.New("exp").Parse(expectedStr))

	tmpl.Execute(&current, js)
	tmplExp.Execute(&expected, tmplExp.Name())

	if current.String() != expected.String() {
		t.Log("Step 1")
		t.Logf("Current is  : %s", current.String())
		t.Logf("Expected is : %s", expected.String())
		t.Fail()
	}

	current.Reset()
	expected.Reset()

	if err := js.IncludeHtml(`<p id="id">hello world</p>`); err != nil {
		t.Errorf("IncludeHtml failed, error : %s", err)
	}

	/*
	 * objForTest={
	 *   _buildStyle:function(){
	 *     var a = document.createElement("style");
	 *     a.innerHTML = "  p {    color: red;  }  #id {    font-size: 2em;  }  ";
	 *     this._doc.appendChild(a);
	 *   },
	 *   _buildDoc:function(target){
	 *     var _d = document;
	 *     var _sr = _d.querySelector(target).attachShadow({mode:'open'});
	 *     var b = _d.createElement("p");
	 *     b.setAttribute("id", "id");
	 *     _sr.appendChild(b);
	 *     var c = _d.createTextNode("hello world");
	 *     b.appendChild(c);
	 *     this._doc = _sr;
	 *   }
	 * }
	 */
	expectedStr = `objForTest={_buildStyle:function(){var a = document.createElement("style");a.innerHTML = "  p {    color: red;  }  #id {    font-size: 2em;  }  ";this._doc().appendChild(a);},_buildDoc:function(target){var _d = document;var _sr = _d.querySelector(target).attachShadow({mode:'open'});var b = _d.createElement("p");b.setAttribute("id", "id");_sr.appendChild(b);var c = _d.createTextNode("hello world");b.appendChild(c);this._doc = function() { return _sr};}}`

	tmpl, err = js.Build()
	if err != nil {
		t.Errorf("Build failed, error : %s", err)
	}
	tmplExp = template.Must(template.New("exp").Parse(expectedStr))

	tmpl.Execute(&current, js)
	tmplExp.Execute(&expected, tmplExp.Name())

	if current.String() != expected.String() {
		t.Logf("Step 2")
		t.Logf("Current is  : %s", current.String())
		t.Logf("Expected is : %s", expected.String())
		t.Fail()
	}
}
