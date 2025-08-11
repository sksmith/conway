package conway

import "fmt"

type TruncateOp struct{}

func (t TruncateOp) Symbol() string {
	return "t"
}

func (t TruncateOp) Name() string {
	return "truncate"
}

func (t TruncateOp) Apply(p *Polyhedron) *Polyhedron {
	trunc := NewPolyhedron("t" + p.Name)

	truncFactor := 1.0 / 3.0

	edgeVertices := make(map[string]*Vertex)

	for _, edge := range p.Edges {
		v1Pos := edge.V1.Position
		v2Pos := edge.V2.Position

		newV1Pos := v1Pos.Add(v2Pos.Sub(v1Pos).Scale(truncFactor))
		newV2Pos := v1Pos.Add(v2Pos.Sub(v1Pos).Scale(1 - truncFactor))

		key1 := edgeVertexKey(edge.ID, edge.V1.ID)
		key2 := edgeVertexKey(edge.ID, edge.V2.ID)

		edgeVertices[key1] = trunc.AddVertex(newV1Pos)
		edgeVertices[key2] = trunc.AddVertex(newV2Pos)
	}

	for _, face := range p.Faces {
		newFaceVertices := allocateVertexSlice(len(face.Vertices) * 2) // Pre-allocate with expected capacity

		for i, vertex := range face.Vertices {
			prevVertex := face.Vertices[(i-1+len(face.Vertices))%len(face.Vertices)]
			nextVertex := face.Vertices[(i+1)%len(face.Vertices)]

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

			if edge1 != nil && edge2 != nil {
				key1 := edgeVertexKey(edge1.ID, vertex.ID)
				key2 := edgeVertexKey(edge2.ID, vertex.ID)

				if v1, ok := edgeVertices[key1]; ok {
					newFaceVertices = append(newFaceVertices, v1)
				}
				if v2, ok := edgeVertices[key2]; ok {
					newFaceVertices = append(newFaceVertices, v2)
				}
			}
		}

		if len(newFaceVertices) >= 3 {
			trunc.AddFace(newFaceVertices)
		}
	}

	for _, vertex := range p.Vertices {
		vertexFaceVertices := allocateVertexSlice(vertex.Degree()) // Pre-allocate based on vertex degree

		orderedEdges := orderEdgesAroundVertex(vertex)
		for _, edge := range orderedEdges {
			key := edgeVertexKey(edge.ID, vertex.ID)
			if v, ok := edgeVertices[key]; ok {
				vertexFaceVertices = append(vertexFaceVertices, v)
			}
		}

		if len(vertexFaceVertices) >= 3 {
			trunc.AddFace(vertexFaceVertices)
		}
	}

	trunc.Normalize()
	return trunc
}

func edgeVertexKey(edgeID, vertexID int) string {
	return fmt.Sprintf("%d_%d", edgeID, vertexID)
}

func Truncate(p *Polyhedron) *Polyhedron {
	op := TruncateOp{}
	return op.Apply(p)
}
