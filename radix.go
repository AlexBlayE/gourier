package gourier

type radixNode struct {
	children map[byte]*radixNode
	handlers []HandleFunc
	depth    uint
}

func (r *radixNode) FindPath(b ...byte) *radixNode {
	size := len(b)
	if size == 0 {
		return nil
	}

	node := r
	for _, elem := range b {
		child, ok := node.children[elem]
		// TODO: aquí mirar si pot tenir mes children o ha arribat al final perque deixi de llegir bytes del payload
		if !ok {
			return nil
		}

		node = child
	}

	return node
}

func (r *radixNode) GetDepth() uint {
	return r.depth
}

func (r *radixNode) GetHandlers() []HandleFunc {
	return r.handlers
}
