package graph

import (
	"reflect"
	"testing"
)

func TestTraverse(t *testing.T) {
	tests := []struct {
		name  string
		graph *Graph
		from  string
		depth int
		want  []string
	}{
		{
			name: "linear chain",
			graph: &Graph{
				Nodes: map[string]Node{"A": {}, "B": {}, "C": {}},
				outgoing: map[string][]Edge{
					"A": {{Target: "B"}},
					"B": {{Target: "C"}},
				},
			},
			from:  "A",
			depth: 10,
			want:  []string{"A", "B", "C"},
		},
		{
			name: "respects depth",
			graph: &Graph{
				Nodes: map[string]Node{"A": {}, "B": {}, "C": {}},
				outgoing: map[string][]Edge{
					"A": {{Target: "B"}},
					"B": {{Target: "C"}},
				},
			},
			from:  "A",
			depth: 1,
			want:  []string{"A", "B"},
		},
		{
			name: "cycle terminates",
			graph: &Graph{
				Nodes: map[string]Node{"A": {}, "B": {}, "C": {}},
				outgoing: map[string][]Edge{
					"A": {{Target: "B"}},
					"B": {{Target: "A"}, {Target: "C"}},
					"C": {{Target: "A"}},
				},
			},
			from:  "A",
			depth: 10,
			want:  []string{"A", "B", "C"},
		},
		{
			name: "single node",
			graph: &Graph{
				Nodes:    map[string]Node{"A": {}},
				outgoing: map[string][]Edge{},
			},
			from:  "A",
			depth: 10,
			want:  []string{"A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.graph.Traverse(tt.from, tt.depth, Forward)
			if err != nil {
				t.Fatalf("Traverse returned error: %v", err)
			}

			ids := make([]string, len(got))
			for i, n := range got {
				ids[i] = n.ID
			}
			if !reflect.DeepEqual(ids, tt.want) {
				t.Errorf("got %v, want %v", ids, tt.want)
			}
		})
	}
}
