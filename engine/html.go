package engine

import (
  h "golang.org/x/net/html"
  "strings"
)

type html struct {
  depthInd  int
  doc       *h.Node
  exclNodes []string
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
    make([]string, 0),
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

func (html *html) readAndExecute(fn func(*html, int)) {
  fn(html, html.depthInd)
  for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
    html.curNode = c
    html.readAndExecute(fn)
  }
  html.depthInd = html.depthInd + 1
}
