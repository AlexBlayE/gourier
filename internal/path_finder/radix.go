package pathfinder

import (
	"github.com/AlexBlayE/gourier"
)

type RadixNode struct {
	Children map[byte]gourier.PathFinder

	Handlers     []gourier.HandleFunc
	ErrorHandler gourier.HandleFunc

	Depth uint
}

func (r *RadixNode) FindPath(b ...byte) gourier.PathFinder {
	size := len(b)
	if size == 0 {
		return nil
	}

	node := r
	for _, elem := range b {
		if node.Children == nil {
			return node
		}

		child, ok := node.Children[elem]
		if !ok {
			return &RadixNode{nil, []gourier.HandleFunc{r.ErrorHandler}, nil, r.Depth + 1}
		}

		node = child.(*RadixNode)
	}

	return node
}

func (r *RadixNode) GetDepth() uint {
	return r.Depth
}

func (r *RadixNode) GetHandlers() []gourier.HandleFunc {
	return r.Handlers
}

func (r *RadixNode) SetErrorHandler(hf gourier.HandleFunc) {
	r.ErrorHandler = hf
}

func (r *RadixNode) SetChild(b byte, child gourier.PathFinder) {
	r.Children[b] = child
}
