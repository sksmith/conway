package conway

type DualOp struct{}

func (d DualOp) Symbol() string {
	return "d"
}

func (d DualOp) Name() string {
	return "dual"
}

func (d DualOp) Apply(p *Polyhedron) *Polyhedron {
	dual := NewPolyhedron("d" + p.Name)

	faceVertices := make(map[int]*Vertex)
	for _, face := range p.Faces {
		centroid := face.Centroid()
		v := dual.AddVertex(centroid)
		faceVertices[face.ID] = v
	}

	for _, edge := range p.Edges {
		if len(edge.Faces) == 2 {
			faces := make([]*Face, 0, 2)
			for _, f := range edge.Faces {
				faces = append(faces, f)
			}

			v1 := faceVertices[faces[0].ID]
			v2 := faceVertices[faces[1].ID]
			dual.AddEdge(v1, v2)
		}
	}

	for _, vertex := range p.Vertices {
		if len(vertex.Faces) >= 3 {
			orderedFaces := orderFacesAroundVertex(vertex)
			dualVertices := make([]*Vertex, len(orderedFaces))
			for i, face := range orderedFaces {
				dualVertices[i] = faceVertices[face.ID]
			}
			dual.AddFace(dualVertices)
		}
	}

	dual.Normalize()
	return dual
}

func orderFacesAroundVertex(v *Vertex) []*Face {
	if len(v.Faces) == 0 {
		return []*Face{}
	}

	faces := make([]*Face, 0, len(v.Faces))
	for _, f := range v.Faces {
		faces = append(faces, f)
	}

	if len(faces) <= 2 {
		return faces
	}

	ordered := make([]*Face, 0, len(faces))
	visited := make(map[int]bool)

	current := faces[0]
	ordered = append(ordered, current)
	visited[current.ID] = true

	for len(ordered) < len(faces) {
		found := false
		for _, edge := range v.Edges {
			for _, face := range edge.Faces {
				if face.ID == current.ID {
					continue
				}

				if visited[face.ID] {
					continue
				}

				hasSharedEdge := false
				for _, e := range current.Edges {
					idx := findEdgeIndex(face, e)
					if idx >= 0 {
						hasSharedEdge = true
						break
					}
				}

				if hasSharedEdge {
					ordered = append(ordered, face)
					visited[face.ID] = true
					current = face
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			for _, f := range faces {
				if !visited[f.ID] {
					ordered = append(ordered, f)
					visited[f.ID] = true
					current = f
					break
				}
			}
		}
	}

	return ordered
}

func findEdgeIndex(face *Face, edge *Edge) int {
	for i, e := range face.Edges {
		if e.ID == edge.ID {
			return i
		}
	}
	return -1
}

func Dual(p *Polyhedron) *Polyhedron {
	op := DualOp{}
	return op.Apply(p)
}
