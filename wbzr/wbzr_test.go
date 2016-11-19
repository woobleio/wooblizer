package wbzr_test

import (
	"strings"
	"testing"

	"github.com/woobleio/wooblizer/wbzr"
)

func TestInject(t *testing.T) {
	wb := wbzr.New(wbzr.JSES5)

	if _, err := wb.Inject("obj={}", "foo"); err != nil {
		t.Errorf("Failed to inject foo, error : %s", err)
	}

	// foo already exists
	if _, err := wb.Inject("otherObj={}", "foo"); err == nil {
		t.Error("It should trigger an error, but it returns a nil error")
	}
}

func TestWrap(t *testing.T) {
	wb := wbzr.New(wbzr.JSES5)

	script1, err := wb.Inject("obj = { _init: function() { console.log('hello'); } }", "obj1")
	if err != nil {
		t.Error("Failed to inject the first script, error : %s", err)
	}

	script2, err := wb.Inject("", "obj2")
	if err != nil {
		t.Error("Failed to inject the second script, error %s", err)
	}

	src := `
	<div id='divid'>
		<p>This is a html and it should be included in the wooble</p>
		<div class='square'></div>
		<div class='square'></div>
		<div class='square'></div>
		<div class='square'></div>
		<div class='square'>
			<p>Hello world!</p>
		</div>
	</div>`

	err = script1.IncludeHtml(src)
	if err != nil {
		t.Error("Failed to include HTML in script1, error : %s", err)
	}

	script2.IncludeHtml("<span></span>")
	script2.IncludeCss("#div { color: red }")

	bf, err := wb.Wrap()
	if err != nil {
		t.Error("Failed to wrap, error %s", err)
	}

	expected := `var cs = {"obj1":{_init:function() { console.log('hello'); },_buildDoc:function(target){var _d = document;var _sr = _d.querySelector(target).attachShadow({mode:'open'});var b = _d.createElement("div");_sr.appendChild(b);this._doc = _sr;}},"obj2":{_buildDoc:function(target){var _d = document;var _sr = _d.querySelector(target).attachShadow({mode:'open'});var b = _d.createElement("span");_sr.appendChild(b);this._doc = _sr;},_buildStyle:function(){var a = document.createElement("style");a.innerHTML = "#div { color: red }";this._doc.appendChild(a);}}}`
	if strings.Contains(bf.String(), expected) {
		t.Logf("expected : %s", expected)
		t.Logf("current : %s", bf.String())
		t.Fail()
	}
}
