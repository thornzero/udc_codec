package udc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TreeNode struct {
	Code     string      `yaml:"code"`
	Title    string      `yaml:"title"`
	Parent   *TreeNode   `yaml:"-"`
	Children []*TreeNode `yaml:"children,omitempty"`
}

func loadTree(filename string) ([]*TreeNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var root []*TreeNode
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}
	return root, nil
}

func buildFlat(nodes []*TreeNode, flat map[string]*TreeNode, parent *TreeNode) {
	for _, n := range nodes {
		n.Parent = parent
		flat[n.Code] = n
		buildFlat(n.Children, flat, n)
	}
}
