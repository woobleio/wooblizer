// Package wbzr provides tool to create a wooble in a given language and to package
// (wrap) some woobles.
package wbzr

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/woobleio/wooblizer/wbzr/engine"
)

// ScriptLang are constants for implemented script languages.
type ScriptLang int

const (
	JSES5 ScriptLang = iota
)

type wbzr struct {
	ext      string
	lang     ScriptLang
	scripts  []engine.Script
	skeleton string
}

// New takes a script language which is used to inject and output a file.
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

// BuildScriptFile make a wooble. The parameter wbName is the name of the injected source
// you want to build
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

// Get returns an injected source.
func (wb *wbzr) Get(name string) (engine.Script, error) {
	for _, sc := range wb.scripts {
		if sc.GetName() == name {
			return sc, nil
		}
	}
	return nil, errors.New("Wooble " + name + " not found")
}

// Inject injects a source code to be wooblized. It takes a name which must be
// unique. Src can be empty, it'll create a default object
func (wb *wbzr) Inject(src string, name string) (engine.Script, error) {
	if _, err := wb.Get(name); err == nil {
		return nil, errors.New("Wooble " + name + " already exists. A name must be unique")
	}
	var sc engine.Script
	switch wb.lang {
	case JSES5:
		sc, err := engine.NewJSES5(src, name)
		if err != nil {
			return nil, err
		}
		wb.scripts = append(wb.scripts, sc)
	}
	return sc, nil
}

// InjectFile injects a source from a file.
func (wb *wbzr) InjectFile(path string, name string) (engine.Script, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return wb.Inject(string(c[:]), name)
}

// Wrap packages some woobles (all the woobles injected in the wbzr)
// and build a file which contains the wooble library.
func (wb *wbzr) Wrap() (*io.Writer, error) {
	for _, sc := range wb.scripts {
		if _, err := sc.Build(); err != nil {
			return nil, err
		}
	}

	tmpl := template.Must(template.New("WbJSES5").Parse(wbJses5))

	var writer io.Writer
	if err := tmpl.Execute(writer, wb.scripts); err != nil {
		return nil, err
	}

	return &writer, nil
}

var wbJses5 = `
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
