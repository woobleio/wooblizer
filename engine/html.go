package engine

import (
  h "golang.org/x/net/html"
  "strings"
)

type html struct {
  exclNodes []string
  doc       *h.Node
  curNode   *h.Node
}

func NewHTML(doc string) (*html, error) {
  r := strings.NewReader(doc)
  node, err := h.Parse(r)
  if err != nil {
    return nil, err
  }

  return &html{
    make([]string, 0),
    node,
    node,
  }, nil
}

func (html *html) AddExcludedNodes(nodes ...string) {
  for _, n := range nodes {
    html.exclNodes = append(html.exclNodes, n)
  }
}

func (html *html) isExcludedNode(node string) bool {
  for _, n := range html.exclNodes {
    if n == node {
      return true
    }
  }
  return false
}

func (html *html) readAndExecute(fn func(*html)) {
  fn(html)
  for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
    html.curNode = c
    html.readAndExecute(fn)
  }
}
