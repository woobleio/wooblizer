package doc

import (
	"regexp"

	"strings"

	h "golang.org/x/net/html"
)

// HTML is struct for HTML parse
type HTML struct {
	doc     *h.Node
	curNode *h.Node
}

var exclNodes []interface{}

// NewHTML creates a new HTML parser
func NewHTML(doc string) (*HTML, error) {
	r := strings.NewReader(doc)
	node, err := h.Parse(r)
	if err != nil {
		return nil, err
	}

	addExcludedNodes("body", "html", "head", h.DoctypeNode, h.ErrorNode, h.DocumentNode, h.CommentNode)

	return &HTML{
		node,
		node,
	}, nil
}

// ReadAndExecute is a recursive function that takes a callback function as parameters.
// It parses the HTML source code and execute a callback function with the current node.
// It's a Depth-first Search algorithm
func (html *HTML) ReadAndExecute(fn func(*h.Node, int) int, pIndex int) {
	n := html.curNode

	// Fixes html string format, avoid " " text nodes, for insecable space use &nbsp;
	isInvalid, _ := regexp.MatchString("^\\s+$", n.Data)
	if isInvalid {
		return
	}

	if !isExcludedNode(n) {
		pIndex = fn(n, pIndex)
	}
	for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
		html.curNode = c
		html.ReadAndExecute(fn, pIndex)
	}
}

// addExcludedNodes add nodes to be excluded from the parser
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
