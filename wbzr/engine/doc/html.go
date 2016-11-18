package doc

import (
	"regexp"

	h "golang.org/x/net/html"
	"strings"
)

type html struct {
	depthInd int
	doc      *h.Node
	curNode  *h.Node
}

var exclNodes []interface{}

func NewHTML(doc string) (*html, error) {
	r := strings.NewReader(doc)
	node, err := h.Parse(r)
	if err != nil {
		return nil, err
	}

	addExcludedNodes("body", "html", "head", h.DoctypeNode, h.ErrorNode, h.DocumentNode, h.CommentNode)

	return &html{
		0,
		node,
		node,
	}, nil
}

func (html *html) ReadAndExecute(fn func(*h.Node, int)) {
	n := html.curNode

	// Fixes html string format, avoid " " text nodes, for insecable space use &nbsp;
	isInvalid, _ := regexp.MatchString("^\\s+$", n.Data)

	if !isExcludedNode(n) && !isInvalid {
		fn(n, html.depthInd)
	}
	for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
		html.curNode = c
		html.ReadAndExecute(fn)
	}
	html.depthInd = html.depthInd + 1
}

func addExcludedNodes(nodes ...interface{}) {
	for _, n := range nodes {
		exclNodes = append(exclNodes, n)
	}
}

func isExcludedNode(node *h.Node) bool {
	for _, n := range exclNodes {
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
	}
	return false
}
