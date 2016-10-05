package script

type Script interface {
  TranspileObject(src string) (*WbObject, error)
}

type attribute struct {
  Type    string
  Value   interface{}
}

type method struct {
  Name    string
  Args    map[string]interface{}
  Return  interface{}
}

type WbObject struct {
  Methods   []method
  Attrs     []attribute
}
