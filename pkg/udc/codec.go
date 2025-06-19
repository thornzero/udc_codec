package udc

import "fmt"

type Codec struct {
	tree *TreeNode
	flat map[string]*TreeNode
}

func LoadCodec(filename string) (*Codec, error) {
	nodes, err := loadTree(filename)
	if err != nil {
		return nil, err
	}
	flat := make(map[string]*TreeNode)
	buildFlat(nodes, flat, nil)

	return &Codec{
		tree: &TreeNode{Children: nodes},
		flat: flat,
	}, nil
}

func (c *Codec) Lookup(code string) (string, bool) {
	node, ok := c.flat[code]
	if !ok {
		return "", false
	}
	return node.Title, true
}

func (c *Codec) ParseComposite(code string) ([]*TreeNode, error) {
	parts, err := parseComposite(code)
	if err != nil {
		return nil, err
	}
	var result []*TreeNode
	for _, p := range parts {
		node, ok := c.flat[p]
		if !ok {
			return nil, fmt.Errorf("unknown code part: %s", p)
		}
		result = append(result, node)
	}
	return result, nil
}

func (c *Codec) Search(term string) []*TreeNode {
	return search(c.flat, term)
}

func (c *Codec) Children(code string) ([]*TreeNode, bool) {
	node, ok := c.flat[code]
	if !ok {
		return nil, false
	}
	return node.Children, true
}

func (c *Codec) Ancestry(code string) ([]*TreeNode, bool) {
	node, ok := c.flat[code]
	if !ok {
		return nil, false
	}
	var path []*TreeNode
	for n := node; n != nil; n = n.Parent {
		path = append([]*TreeNode{n}, path...)
	}
	return path, true
}

func (c *Codec) Validate(code string) error {
	return validateComposite(code, c.flat)
}
