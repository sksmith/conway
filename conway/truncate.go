package conway

import "fmt"

const (
	// defaultTruncateFactor is the standard truncation factor (1/3).
	defaultTruncateFactor = 1.0 / 3.0
)

type TruncateOp struct{}

func (t TruncateOp) Symbol() string {
	return "t"
}

func (t TruncateOp) Name() string {
	return "truncate"
}

// createTruncatedEdgeVertices creates new vertices along each edge for truncation.
func createTruncatedEdgeVertices(p *Polyhedron, truncFactor float64) (
	map[string]*Vertex, *Polyhedron,
) {
	trunc := NewPolyhedron("t" + p.Name)

	edgeVertices := make(map[string]*Vertex)

	for _, edge := range p.Edges {
		v1Pos := edge.V1.Position

		v2Pos := edge.V2.Position

		newV1Pos := v1Pos.Add(v2Pos.Sub(v1Pos).Scale(truncFactor))

		newV2Pos := v1Pos.Add(v2Pos.Sub(v1Pos).Scale(1 - truncFactor))

		key1 := EdgeVertexKey(edge.ID, edge.V1.ID)

		key2 := EdgeVertexKey(edge.ID, edge.V2.ID)

		edgeVertices[key1] = trunc.AddVertex(newV1Pos)
		edgeVertices[key2] = trunc.AddVertex(newV2Pos)
	}

	return edgeVertices, trunc
}

// findAdjacentEdges finds the edges connecting a vertex to its previous and next neighbors in a face.
func findAdjacentEdges(vertex, prevVertex, nextVertex *Vertex) (*Edge, *Edge) {
	var edge1, edge2 *Edge

	for _, e := range vertex.Edges {
		other := e.OtherVertex(vertex)

		if other != nil {
			if other.ID == prevVertex.ID {
				edge1 = e
			} else if other.ID == nextVertex.ID {
				edge2 = e
			}
		}
	}

	return edge1, edge2
}

// addTruncatedFaceVertices adds vertices for a truncated face.
func addTruncatedFaceVertices(face *Face, edgeVertices map[string]*Vertex) []*Vertex {
	newFaceVertices := allocateVertexSlice(len(face.Vertices) * 2)

	for i, vertex := range face.Vertices {
		prevVertex := face.Vertices[(i-1+len(face.Vertices))%len(face.Vertices)]

		nextVertex := face.Vertices[(i+1)%len(face.Vertices)]

		edge1, edge2 := findAdjacentEdges(vertex, prevVertex, nextVertex)

		if edge1 != nil && edge2 != nil {
			key1 := EdgeVertexKey(edge1.ID, vertex.ID)

			key2 := EdgeVertexKey(edge2.ID, vertex.ID)

			if v1, ok := edgeVertices[key1]; ok {
				newFaceVertices = append(newFaceVertices, v1)
			}

			if v2, ok := edgeVertices[key2]; ok {
				newFaceVertices = append(newFaceVertices, v2)
			}
		}
	}

	return newFaceVertices
}

// processTruncatedFaces processes all faces to create truncated versions.
func processTruncatedFaces(p, trunc *Polyhedron, edgeVertices map[string]*Vertex) {
	for _, face := range p.Faces {
		newFaceVertices := addTruncatedFaceVertices(face, edgeVertices)

		if len(newFaceVertices) >= 3 {
			trunc.AddFace(newFaceVertices)
		}
	}
}

// processTruncatedVertexFaces processes vertices to create new faces at truncation sites.
func processTruncatedVertexFaces(p, trunc *Polyhedron, edgeVertices map[string]*Vertex) {
	for _, vertex := range p.Vertices {
		vertexFaceVertices := allocateVertexSlice(vertex.Degree())

		orderedEdges := OrderEdgesAroundVertex(vertex)

		for _, edge := range orderedEdges {
			key := EdgeVertexKey(edge.ID, vertex.ID)

			if v, ok := edgeVertices[key]; ok {
				vertexFaceVertices = append(vertexFaceVertices, v)
			}
		}

		if len(vertexFaceVertices) >= 3 {
			trunc.AddFace(vertexFaceVertices)
		}
	}
}

func (t TruncateOp) Apply(p *Polyhedron) *Polyhedron {
	edgeVertices, trunc := createTruncatedEdgeVertices(p, defaultTruncateFactor)
	processTruncatedFaces(p, trunc, edgeVertices)
	processTruncatedVertexFaces(p, trunc, edgeVertices)
	trunc.Normalize()

	return trunc
}

func EdgeVertexKey(edgeID, vertexID int) string {
	return fmt.Sprintf("%d_%d", edgeID, vertexID)
}

func Truncate(p *Polyhedron) *Polyhedron {
	op := TruncateOp{}
	return op.Apply(p)
}
