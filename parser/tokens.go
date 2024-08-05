package parser

type TokenGraphNode[T any] struct {
	children map[rune]*TokenGraphNode[T]
	value    T
}

func createTokenGraph[D any](stringMap *map[string]D) *TokenGraphNode[D] {
	rootNode := TokenGraphNode[D]{
		children: make(map[rune]*TokenGraphNode[D]),
	}

	for key, val := range *stringMap {
		node := &rootNode
		for _, token := range key {
			childNode, hasToken := node.children[token]
			if !hasToken {
				childNode = &TokenGraphNode[D]{
					children: make(map[rune]*TokenGraphNode[D]),
				}
				node.children[token] = childNode
			}
			node = childNode
		}
		node.value = val
	}

	return &rootNode
}
