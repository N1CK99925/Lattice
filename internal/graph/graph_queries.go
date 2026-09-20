package graph

type Item struct {
	id    string
	depth int
}

type Direction int

const (
	Forward Direction = iota
	Reverse
)

// BFS with depth for traversal to get dependents, dependecies and deadcode
func (g *Graph) Traverse(from string, Depth int, direction Direction) ([]Node, error) {
	visited := map[string]bool{from: true}
	queue := []Item{{id: from, depth: 0}}
	var order []string
	var Dir map[string][]Edge

	if direction == Reverse {
		Dir = g.incoming

	} else if direction == Forward {
		Dir = g.outgoing
	}
	for len(queue) > 0 {

		current := queue[0]
		queue = queue[1:]
		order = append(order, current.id)
		for _, edge := range Dir[current.id] {

			// var sex suggested by Rishabh
			var sex string
			if direction == Reverse {
				sex = edge.Source
			} else if direction == Forward {
				sex = edge.Target
			}
			nextDepth := current.depth + 1
			if nextDepth > Depth {
				continue
			}
			if !visited[sex] {
				visited[sex] = true
				queue = append(queue, Item{
					id:    sex,
					depth: nextDepth,
				})

			}
		}

	}
	var nodes []Node
	for _, id := range order {
		node := Node{ID: id}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Things that A depend on
func (g *Graph) dependecies(symbol string, depth int) ([]Node, error) {

	return g.Traverse(symbol, depth, Forward)
}

// Thing that depend on A
func (g *Graph) dependents(symbol string, depth int) ([]Node, error) {
	return g.Traverse(symbol, depth, Forward)
}
