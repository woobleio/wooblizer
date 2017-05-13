package wbzr_test

import (
	"testing"

	"github.com/woobleio/wooblizer/wbzr"
)

func TestInject(t *testing.T) {
	wb := wbzr.New(wbzr.JS)
	if _, err := wb.Inject("obj={}", "foo"); err != nil {
		t.Errorf("Failed to inject foo, error : %s", err)
	}

	// foo already exists
	if _, err := wb.Inject("otherObj={}", "foo"); err == nil {
		t.Error("It should trigger an error, but it returns a nil error")
	}
}

func TestSecureAndWrap(t *testing.T) {
	wb := wbzr.New(wbzr.JS)

	/*script1, err := wb.Inject(`var Woobly = function () {
	  function Woobly(toto) {
	    _classCallCheck(this, Woobly);
			this.document = document;
	  }

	  _createClass(Woobly, [{
	    key: "getDocument",
	    value: function getDocument() {
	      return document;
	    }
	  }]);

	  return Woobly;
	}();`, "obj1")

		if err != nil {
			t.Error("Failed to inject the first script, error : %s", err)
		}*/

	script2, _ := wb.Inject(`var Woobly = function(){function Woobly(){_classCallCheck(this,Woobly);this.document=document}_createClass(Woobly,[{key:"toto",value:function toto(lol){}}]);return Woobly}();`, "obj2")

	// src := `
	// <div id='divid'>
	// 	yoyo
	// </div>`

	// err = script1.IncludeHTMLCSS(src, "div { color: red; }")
	// if err != nil {
	// 	t.Error("Failed to include HTML in script1, error : %s", err)
	// }

	err := script2.IncludeHTMLCSS("", "div { color: red; }")
	if err != nil {
		t.Error("Failed to include HTML in script2, error : %s", err)
	}

	bf, err := wb.SecureAndWrap("toto.com", "tata.com")
	if err != nil {
		t.Error("Failed to wrap, error %s", err)
	}

	t.Log(bf.String())
}
