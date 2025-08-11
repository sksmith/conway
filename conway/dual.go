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
		if len(edge.Faces) != 2 {
			continue
		}

		faces := make([]*Face, 0, 2)

		for _, f := range edge.Faces {
			faces = append(faces, f)
		}

		v1 := faceVertices[faces[0].ID]

		v2 := faceVertices[faces[1].ID]

		dual.AddEdge(v1, v2)
	}

	for _, vertex := range p.Vertices {
		if len(vertex.Faces) >= 3 {
			orderedFaces := OrderFacesAroundVertex(vertex)

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

// convertFacesToSlice converts vertex faces map to slice.
func convertFacesToSlice(v *Vertex) []*Face {
	faces := make([]*Face, 0, len(v.Faces))

	for _, f := range v.Faces {
		faces = append(faces, f)
	}

	return faces
}

// facesShareEdge checks if two faces share an edge.
func facesShareEdge(face1, face2 *Face) bool {
	for _, e := range face1.Edges {
		if FindEdgeIndex(face2, e) >= 0 {
			return true
		}
	}

	return false
}

// findNextFaceInEdges searches through vertex edges to find the next adjacent face.
func findNextFaceInEdges(v *Vertex, currentFace *Face, visited map[int]bool) *Face {
	for _, edge := range v.Edges {
		for _, face := range edge.Faces {
			if face.ID != currentFace.ID && !visited[face.ID] {
				if facesShareEdge(currentFace, face) {
					return face
				}
			}
		}
	}

	return nil
}

// findNextUnvisitedFace finds any unvisited face (fallback).
func findNextUnvisitedFace(faces []*Face, visited map[int]bool) *Face {
	for _, f := range faces {
		if !visited[f.ID] {
			return f
		}
	}

	return nil
}

func OrderFacesAroundVertex(v *Vertex) []*Face {
	if len(v.Faces) == 0 {
		return []*Face{}
	}

	faces := convertFacesToSlice(v)

	if len(faces) <= 2 {
		return faces
	}

	ordered := make([]*Face, 0, len(faces))

	visited := make(map[int]bool)

	current := faces[0]

	ordered = append(ordered, current)
	visited[current.ID] = true

	for len(ordered) < len(faces) {
		// Try to find next face through edge connections.
		nextFace := findNextFaceInEdges(v, current, visited)

		// If no edge-connected face found, use fallback.
		if nextFace == nil {
			nextFace = findNextUnvisitedFace(faces, visited)
		}

		if nextFace != nil {
			ordered = append(ordered, nextFace)
			visited[nextFace.ID] = true
			current = nextFace
		} else {
			break // Safety break in case we can't find any more faces
		}
	}

	return ordered
}

func FindEdgeIndex(face *Face, edge *Edge) int {
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
