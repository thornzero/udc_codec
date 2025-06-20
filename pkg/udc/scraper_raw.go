package udc


import (
	"context"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

type RawNode struct {
	ID       string
	Parent   string
	Code     string
	Title    string
	Children []*RawNode
}

func ScrapeRawTree(ctx context.Context) ([]*RawNode, error) {
	var rawHTML string
	err := chromedp.Run(ctx,
		chromedp.OuterHTML("#classtree", &rawHTML),
	)
	if err != nil {
		return nil, err
	}

	parsed := parseRawHTML(rawHTML)
	root := buildRawHierarchy(parsed)
	return root, nil
}

func parseRawHTML(html string) []*RawNode {
	re := regexp.MustCompile(`d\.add\((\d+),(\d+),'(.*?)','.*?&nbsp;&nbsp;(.*?)','`)
	matches := re.FindAllStringSubmatch(html, -1)

	var nodes []*RawNode
	for _, m := range matches {
		code := strings.TrimSpace(m[3])
		if code == "-" || code == "--" || code == "---" || code == "----" {
			continue
		}
		node := &RawNode{
			ID:     m[1],
			Parent: m[2],
			Code:   code,
			Title:  strings.TrimSpace(m[4]),
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func buildRawHierarchy(flat []*RawNode) []*RawNode {
	idMap := make(map[string]*RawNode)
	for _, n := range flat {
		idMap[n.ID] = n
	}
	var roots []*RawNode
	for _, n := range flat {
		if n.Parent == "0" || n.ID == "0" {
			roots = append(roots, n)
			continue
		}
		if parent, ok := idMap[n.Parent]; ok {
			parent.Children = append(parent.Children, n)
		}
	}
	return roots
}
