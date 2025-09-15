package group

import (
	"github.com/AlexBlayE/gourier"
	pathfinder "github.com/AlexBlayE/gourier/internal/path_finder"
)

type RouterGroup struct {
	PathFinder gourier.PathFinder
}

func (rg *RouterGroup) Group(header byte) gourier.RouteGroup {
	newChildNode := &pathfinder.RadixNode{make(map[byte]gourier.PathFinder), nil, nil, rg.PathFinder.GetDepth() + 1}
	rg.PathFinder.SetChild(header, newChildNode)
	return &RouterGroup{newChildNode}
}

// TODO: fer que es tingui que injectar el contructor de pathfinder
func (rg *RouterGroup) Handler(header byte, handleFunc ...gourier.HandleFunc) {
	node := &pathfinder.RadixNode{nil, handleFunc, nil, rg.PathFinder.GetDepth() + 1}
	rg.PathFinder.SetChild(header, node)
}

func (rg *RouterGroup) Error(errorHandler gourier.HandleFunc) {
	rg.PathFinder.SetErrorHandler(errorHandler)
}
