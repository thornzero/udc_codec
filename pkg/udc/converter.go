package udc

func ConvertRawToModel(raw []*RawNode) []*UDCNode {
	var convert func(r *RawNode) *UDCNode
	convert = func(r *RawNode) *UDCNode {
		node := &UDCNode{
			Code:  r.Code,
			Title: r.Title,
		}
		for _, child := range r.Children {
			node.Children = append(node.Children, convert(child))
		}
		return node
	}
	var roots []*UDCNode
	for _, r := range raw {
		roots = append(roots, convert(r))
	}
	return roots
}
