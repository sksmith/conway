package conway

type KisOp struct{}

func (k KisOp) Symbol() string {
	return "k"
}

func (k KisOp) Name() string {
	return "kis"
}

func (k KisOp) Apply(p *Polyhedron) *Polyhedron {
	kis := NewPolyhedron("k" + p.Name)

	vertexMap := make(map[int]*Vertex)
	for _, v := range p.Vertices {
		newV := kis.AddVertex(v.Position)
		vertexMap[v.ID] = newV
	}

	pyramidHeight := 0.5

	for _, face := range p.Faces {
		centroid := face.Centroid()
		normal := face.Normal()

		apexPos := centroid.Add(normal.Scale(pyramidHeight))
		apex := kis.AddVertex(apexPos)

		faceVertices := make([]*Vertex, len(face.Vertices))
		for i, v := range face.Vertices {
			faceVertices[i] = vertexMap[v.ID]
		}

		for i := 0; i < len(faceVertices); i++ {
			v1 := faceVertices[i]
			v2 := faceVertices[(i+1)%len(faceVertices)]
			kis.AddFace([]*Vertex{v1, v2, apex})
		}
	}

	kis.Normalize()
	return kis
}

func Kis(p *Polyhedron) *Polyhedron {
	op := KisOp{}
	return op.Apply(p)
}
