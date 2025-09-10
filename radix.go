package gourier

type radixNode struct {
	children map[byte]*radixNode

	handlers     []HandleFunc
	errorHandler HandleFunc // TODO:

	depth uint
}

func (r *radixNode) FindPath(b ...byte) *radixNode {
	size := len(b)
	if size == 0 {
		return nil
	}

	node := r
	for _, elem := range b {
		if node.children == nil {
			return node
		}

		child, ok := node.children[elem]
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
