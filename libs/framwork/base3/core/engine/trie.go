package engine

type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的一部分
	children []*node // 子节点
	isWild   bool    // 是否精确匹配。part 含有 * 或 : 为 true
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
