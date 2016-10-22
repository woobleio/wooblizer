package wbzr

import (
  "errors"
  "os"
  "text/template"

  "github.com/woobleio/wooblizer/wbzr/engine"
)

type ScriptLang int
const (
  JSES5 ScriptLang = iota
)

type wbzr struct {
  ext       string
  lang      ScriptLang
  scripts   []engine.Script
  skeleton  string
}

func New(sl ScriptLang) *wbzr {
  var skeleton, ext string
  switch sl {
  case JSES5:
    skeleton = wbJses5
    ext = ".min.js"
  default:
    panic("Language not supported")
  }

  return &wbzr{
    ext,
    sl,
    make([]engine.Script, 0),
    skeleton,
  }
}

func (wb *wbzr) BuildScriptFile(wbName string, path string, fileName string) error {
  script, err := wb.Get(wbName)
  if err != nil {
    return err
  }

  tmplSrc, err := script.Build()
  if err != nil {
    return err
  }

  f, err := os.Create(path + "/" + fileName + wb.ext)
  if err != nil {
    return err
  }
  defer f.Close()

  if err := tmplSrc.Execute(f, script); err != nil {
    return err
  }

  return nil
}

func (wb *wbzr) Get(name string) (engine.Script, error) {
  for _, sc := range wb.scripts {
    if sc.GetName() == name {
      return sc, nil
    }
  }
  return nil, errors.New("Wooble " + name + " not found")
}

func (wb *wbzr) Inject(src string, name string) (*engine.Script, error) {
  var sc engine.Script
  switch wb.lang {
  case JSES5:
    sc, err := engine.NewJSES5(src, name)
    if err != nil {
      return nil, err
    }
    wb.scripts = append(wb.scripts, sc)
  }
  return &sc, nil
}

func (wb *wbzr) WrapAndBuildFile(path string, fileName string) error {
  for _, sc := range wb.scripts {
    if _, err := sc.Build(); err != nil {
      return err
    }
  }

  f, err := os.Create(path + "/" + fileName + wb.ext)
  if err != nil {
    return err
  }
  defer f.Close()

  tmpl := template.Must(template.New("WbJSES5").Parse(wbJses5))

  if err := tmpl.Execute(f, wb.scripts); err != nil {
    return err
  }

  return nil
}

var wbJses5 =
`
function Wb(id) {
	if(window === this) {
  	return new Wb(id);
  }

  var cs = {
  	{{range $i, $o := .}}"{{$o.GetName}}":{{"{"}}{{$o.GetSource}}{{"}"}}{{if not $i}},{{end}}{{end}}
  }

  var c = cs[id];
  if(typeof c == 'undefined') {
  	console.log("creation", id, "not found");
    return undefined;
  }

  this.init = function (target) {
    if(document.querySelector(target) == null) {
    	console.log("Element", target, "not found in the document");
      return;
    }

    if("_buildDoc" in c) c._buildDoc(target);
    if("_buildStyle" in c) c._buildStyle();
    if("_init" in c) c._init();
  }

  this.get = function() {
  	return c;
  }

  return this;
}
`
