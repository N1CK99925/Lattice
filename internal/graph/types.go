package graph

type Graph struct {
	Nodes map[string]Node
	Edges []Edge

	incoming map[string][]Edge
	outgoing map[string][]Edge
}

type Node struct {
	ID string
}

type Edge struct {
	Source string
	Target string
	Kind   string
}
