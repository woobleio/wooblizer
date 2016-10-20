package doc

import (
  h "golang.org/x/net/html"
  "strings"
)

type html struct {
  depthInd  int
  doc       *h.Node
  exclNodes []interface{}
  curNode   *h.Node
}

func NewHTML(doc string) (*html, error) {
  r := strings.NewReader(doc)
  node, err := h.Parse(r)
  if err != nil {
    return nil, err
  }

  return &html{
    0,
    node,
    make([]interface{}, 0),
    node,
  }, nil
}

func (html *html) AddExcludedNodes(nodes ...interface{}) {
  for _, n := range nodes {
    html.exclNodes = append(html.exclNodes, n)
  }
}

func (html *html) ReadAndExecute(fn func(*h.Node, int)) {
  n := html.curNode
  if !html.isExcludedNode(n) {
    fn(n, html.depthInd)
  }
  for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
    html.curNode = c
    html.ReadAndExecute(fn)
  }
  html.depthInd = html.depthInd + 1
}


func (html *html) isExcludedNode(node *h.Node) bool {
  for _, n := range html.exclNodes {
    switch n.(type) {
    case string:
      if node.Data == n {
        return true
      }
    case h.NodeType:
      if node.Type == n {
        return true
      }
    }
    if node.Data == n || node.Type == n {
      return true
    }
  }
  return false
}
