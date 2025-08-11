package conway

type AmboOp struct{}

func (a AmboOp) Symbol() string {
	return "a"
}

func (a AmboOp) Name() string {
	return "ambo"
}

func (a AmboOp) Apply(p *Polyhedron) *Polyhedron {
	ambo := NewPolyhedron("a" + p.Name)

	edgeVertices := make(map[int]*Vertex)
	for _, edge := range p.Edges {
		midpoint := edge.Midpoint()
		v := ambo.AddVertex(midpoint)
		edgeVertices[edge.ID] = v
	}

	for _, face := range p.Faces {
		faceVertices := make([]*Vertex, len(face.Edges))
		for i, edge := range face.Edges {
			faceVertices[i] = edgeVertices[edge.ID]
		}
		ambo.AddFace(faceVertices)
	}

	for _, vertex := range p.Vertices {
		if len(vertex.Edges) >= 3 {
			orderedEdges := orderEdgesAroundVertex(vertex)
			vertexFaceVertices := make([]*Vertex, len(orderedEdges))
			for i, edge := range orderedEdges {
				vertexFaceVertices[i] = edgeVertices[edge.ID]
			}
			ambo.AddFace(vertexFaceVertices)
		}
	}

	ambo.Normalize()
	return ambo
}

func orderEdgesAroundVertex(v *Vertex) []*Edge {
	if len(v.Edges) == 0 {
		return []*Edge{}
	}

	edges := make([]*Edge, 0, len(v.Edges))
	for _, e := range v.Edges {
		edges = append(edges, e)
	}

	if len(edges) <= 2 {
		return edges
	}

	ordered := make([]*Edge, 0, len(edges))
	visited := make(map[int]bool)

	current := edges[0]
	ordered = append(ordered, current)
	visited[current.ID] = true

	for len(ordered) < len(edges) {
		found := false

		for _, face := range v.Faces {
			hasCurrentEdge := false
			for _, e := range face.Edges {
				if e.ID == current.ID {
					hasCurrentEdge = true
					break
				}
			}

			if !hasCurrentEdge {
				continue
			}

			for _, e := range face.Edges {
				if e.ID == current.ID || visited[e.ID] {
					continue
				}

				if e.V1.ID == v.ID || e.V2.ID == v.ID {
					ordered = append(ordered, e)
					visited[e.ID] = true
					current = e
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if !found {
			for _, e := range edges {
				if !visited[e.ID] {
					ordered = append(ordered, e)
					visited[e.ID] = true
					current = e
					break
				}
			}
		}
	}

	return ordered
}

func Ambo(p *Polyhedron) *Polyhedron {
	op := AmboOp{}
	return op.Apply(p)
}
