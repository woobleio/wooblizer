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

func TestSecureAndWrap(t *testing.T) {
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
	</div>`

	err = script1.IncludeHtml(src)
	if err != nil {
		t.Error("Failed to include HTML in script1, error : %s", err)
	}

	script2.IncludeHtml("<span></span>")
	script2.IncludeCss("#div { color: red }")

	bf, err := wb.SecureAndWrap("toto.com", "tata.com")
	if err != nil {
		t.Error("Failed to wrap, error %s", err)
	}

	expected := `var ah = ["toto.com","tata.com"];var xx = ah.indexOf(window.location.hostname);if(ah.indexOf(window.location.hostname) == -1) {console.log("Wooble error : domain restricted");return;}`
	if strings.Contains(bf.String(), expected) {
		t.Logf("expected : %s", expected)
		t.Logf("current : %s", bf.String())
		t.Fail()
	}
}
