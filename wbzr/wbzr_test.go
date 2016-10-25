package wbzr_test

import (
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
