package main

import (
  "github.com/woobleio/wooblizer/wbzr"
  //"github.com/robertkrimen/otto"
  "log"
)

func main() {
  src :=
  `
  obj = {
    varBool: true,
    varNumFloat: -10.5,
    varNumInt: 10,
    varString: "Hello W3rld!",
    varArray: [1, 2, 3, -1],
    varArrayStr: ["str1", "str2"]
    varObj: {
      a: "It's a string",
      b: false,
      subObj: {
        c: 12,
        d: 0.2,
        e: [1, 2, 3],
        f: { field: "hi" }
      }
    },
    fnWithoutArgs: function() {
      var a = 1;
      var b = 2;
      var c = a + b - (a * b);
      console.log(c, "goog")
    },
    fnWithArgs: function(a, b) {
      console.log(a - b);
    }
  };
  `

  wbzr := wbzr.New(wbzr.JS, src, "bonjour")
  wbzr.BuildFile("/tmp", "aurevoir")

  log.Print("---- END ----")

}
