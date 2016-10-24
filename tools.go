package main

import (
  "github.com/woobleio/wooblizer/wbzr"
)

var js1 =
`
obj1 = {
  id: "yo",
  _init: function() {
    document.querySelector(this.id);
    document.querySelectorAll(this.id);
  }
}
`

var js2 =
`
obj2 = {
  _init: function() {
    document.querySelectorAll("div");
  }
}
`

func main() {
  // FIXME temp for test only
  wb := wbzr.New(wbzr.JSES5)
  wb.Inject(js1, "firstObj")
  wb.Inject(js2, "secObj")
  wb.BuildScriptFile("firstObj", "/tmp", "first_script")

  sc, _ := wb.Get("firstObj")
  sc.IncludeHtml("<div id='test'>hello</div>");
  sc.IncludeCss("#test{ background-color: red; }")

  wb.BuildScriptFile("firstObj", "/tmp", "first_script_bis")

  wb.WrapAndBuildFile("/tmp", "wrap")
}
