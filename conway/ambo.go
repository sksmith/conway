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

// convertEdgesToSlice converts vertex edges map to slice
func convertEdgesToSlice(v *Vertex) []*Edge {
	edges := make([]*Edge, 0, len(v.Edges))
	for _, e := range v.Edges {
		edges = append(edges, e)
	}

	return edges
}

// faceContainsEdge checks if a face contains the given edge
func faceContainsEdge(face *Face, edgeID int) bool {
	for _, e := range face.Edges {
		if e.ID == edgeID {
			return true
		}
	}

	return false
}

// edgeConnectsToVertex checks if an edge connects to the given vertex
func edgeConnectsToVertex(edge *Edge, vertexID int) bool {
	return edge.V1.ID == vertexID || edge.V2.ID == vertexID
}

// findNextEdgeInFace finds the next unvisited edge in a face that connects to vertex
func findNextEdgeInFace(face *Face, currentEdgeID, vertexID int, visited map[int]bool) *Edge {
	for _, e := range face.Edges {
		if e.ID == currentEdgeID || visited[e.ID] {
			continue
		}
		if edgeConnectsToVertex(e, vertexID) {
			return e
		}
	}

	return nil
}

// findNextEdgeInFaces searches through faces to find the next edge to add
func findNextEdgeInFaces(v *Vertex, currentEdge *Edge, visited map[int]bool) *Edge {
	for _, face := range v.Faces {
		if !faceContainsEdge(face, currentEdge.ID) {
			continue
		}

		if nextEdge := findNextEdgeInFace(face, currentEdge.ID, v.ID, visited); nextEdge != nil {
			return nextEdge
		}
	}

	return nil
}

// findNextUnvisitedEdge finds any unvisited edge (fallback)
func findNextUnvisitedEdge(edges []*Edge, visited map[int]bool) *Edge {
	for _, e := range edges {
		if !visited[e.ID] {
			return e
		}
	}

	return nil
}

func orderEdgesAroundVertex(v *Vertex) []*Edge {
	if len(v.Edges) == 0 {
		return []*Edge{}
	}

	edges := convertEdgesToSlice(v)
	if len(edges) <= 2 {
		return edges
	}

	ordered := make([]*Edge, 0, len(edges))
	visited := make(map[int]bool)

	current := edges[0]
	ordered = append(ordered, current)
	visited[current.ID] = true

	for len(ordered) < len(edges) {
		// Try to find next edge through face connections
		nextEdge := findNextEdgeInFaces(v, current, visited)

		// If no face-connected edge found, use fallback
		if nextEdge == nil {
			nextEdge = findNextUnvisitedEdge(edges, visited)
		}

		if nextEdge != nil {
			ordered = append(ordered, nextEdge)
			visited[nextEdge.ID] = true
			current = nextEdge
		} else {
			break // Safety break in case we can't find any more edges
		}
	}

	return ordered
}

func Ambo(p *Polyhedron) *Polyhedron {
	op := AmboOp{}
	return op.Apply(p)
}
