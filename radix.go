package gourier

type radixNode struct {
	children map[byte]*radixNode

	handlers     []HandleFunc
	errorHandler HandleFunc

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
			return &radixNode{nil, []HandleFunc{r.errorHandler}, nil, r.depth + 1} // If path dont'exist return errorHandler
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
