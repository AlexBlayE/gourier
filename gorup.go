package gourier

type routerGroup struct {
	radixNode *radixNode
}

func (rg *routerGroup) Group(header byte) *routerGroup {
	newChildNode := &radixNode{make(map[byte]*radixNode), nil, rg.radixNode.depth + 1}
	rg.radixNode.children[header] = newChildNode
	return &routerGroup{newChildNode}
}

func (rg *routerGroup) Handler(header byte, handleFunc ...HandleFunc) {
	rg.radixNode.children[header] = &radixNode{nil, handleFunc, rg.radixNode.depth + 1}
}
