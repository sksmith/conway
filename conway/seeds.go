package conway

import "math"

const (
	// goldenRatioBase is the square root of 5 used in golden ratio calculation
	goldenRatioBase = 5
	// goldenRatioDivisor is the divisor used in golden ratio calculation
	goldenRatioDivisor = 2.0
)

func Tetrahedron() *Polyhedron {
	p := NewPolyhedron("Tetrahedron")

	a := 1.0 / math.Sqrt(3)
	vertices := []*Vertex{
		p.AddVertex(Vector3{a, a, a}),
		p.AddVertex(Vector3{a, -a, -a}),
		p.AddVertex(Vector3{-a, a, -a}),
		p.AddVertex(Vector3{-a, -a, a}),
	}

	p.AddFace([]*Vertex{vertices[0], vertices[1], vertices[2]})
	p.AddFace([]*Vertex{vertices[0], vertices[1], vertices[3]})
	p.AddFace([]*Vertex{vertices[0], vertices[2], vertices[3]})
	p.AddFace([]*Vertex{vertices[1], vertices[2], vertices[3]})

	p.Normalize()

	return p
}

func Cube() *Polyhedron {
	p := NewPolyhedron("Cube")

	vertices := []*Vertex{
		p.AddVertex(Vector3{1, 1, 1}),
		p.AddVertex(Vector3{1, 1, -1}),
		p.AddVertex(Vector3{1, -1, 1}),
		p.AddVertex(Vector3{1, -1, -1}),
		p.AddVertex(Vector3{-1, 1, 1}),
		p.AddVertex(Vector3{-1, 1, -1}),
		p.AddVertex(Vector3{-1, -1, 1}),
		p.AddVertex(Vector3{-1, -1, -1}),
	}

	p.AddFace([]*Vertex{vertices[0], vertices[2], vertices[3], vertices[1]})
	p.AddFace([]*Vertex{vertices[4], vertices[5], vertices[7], vertices[6]})
	p.AddFace([]*Vertex{vertices[0], vertices[1], vertices[5], vertices[4]})
	p.AddFace([]*Vertex{vertices[2], vertices[6], vertices[7], vertices[3]})
	p.AddFace([]*Vertex{vertices[0], vertices[4], vertices[6], vertices[2]})
	p.AddFace([]*Vertex{vertices[1], vertices[3], vertices[7], vertices[5]})

	p.Normalize()

	return p
}

func Octahedron() *Polyhedron {
	p := NewPolyhedron("Octahedron")

	vertices := []*Vertex{
		p.AddVertex(Vector3{1, 0, 0}),
		p.AddVertex(Vector3{-1, 0, 0}),
		p.AddVertex(Vector3{0, 1, 0}),
		p.AddVertex(Vector3{0, -1, 0}),
		p.AddVertex(Vector3{0, 0, 1}),
		p.AddVertex(Vector3{0, 0, -1}),
	}

	p.AddFace([]*Vertex{vertices[0], vertices[2], vertices[4]})
	p.AddFace([]*Vertex{vertices[0], vertices[4], vertices[3]})
	p.AddFace([]*Vertex{vertices[0], vertices[3], vertices[5]})
	p.AddFace([]*Vertex{vertices[0], vertices[5], vertices[2]})
	p.AddFace([]*Vertex{vertices[1], vertices[4], vertices[2]})
	p.AddFace([]*Vertex{vertices[1], vertices[3], vertices[4]})
	p.AddFace([]*Vertex{vertices[1], vertices[5], vertices[3]})
	p.AddFace([]*Vertex{vertices[1], vertices[2], vertices[5]})

	p.Normalize()

	return p
}

func Dodecahedron() *Polyhedron {
	p := NewPolyhedron("Dodecahedron")

	phi := (1.0 + math.Sqrt(goldenRatioBase)) / goldenRatioDivisor
	invPhi := 1.0 / phi

	vertices := []*Vertex{
		p.AddVertex(Vector3{1, 1, 1}),
		p.AddVertex(Vector3{1, 1, -1}),
		p.AddVertex(Vector3{1, -1, 1}),
		p.AddVertex(Vector3{1, -1, -1}),
		p.AddVertex(Vector3{-1, 1, 1}),
		p.AddVertex(Vector3{-1, 1, -1}),
		p.AddVertex(Vector3{-1, -1, 1}),
		p.AddVertex(Vector3{-1, -1, -1}),

		p.AddVertex(Vector3{0, phi, invPhi}),
		p.AddVertex(Vector3{0, phi, -invPhi}),
		p.AddVertex(Vector3{0, -phi, invPhi}),
		p.AddVertex(Vector3{0, -phi, -invPhi}),

		p.AddVertex(Vector3{invPhi, 0, phi}),
		p.AddVertex(Vector3{invPhi, 0, -phi}),
		p.AddVertex(Vector3{-invPhi, 0, phi}),
		p.AddVertex(Vector3{-invPhi, 0, -phi}),

		p.AddVertex(Vector3{phi, invPhi, 0}),
		p.AddVertex(Vector3{phi, -invPhi, 0}),
		p.AddVertex(Vector3{-phi, invPhi, 0}),
		p.AddVertex(Vector3{-phi, -invPhi, 0}),
	}

	p.AddFace([]*Vertex{vertices[0], vertices[8], vertices[4], vertices[14], vertices[12]})
	p.AddFace([]*Vertex{vertices[0], vertices[12], vertices[2], vertices[17], vertices[16]})
	p.AddFace([]*Vertex{vertices[0], vertices[16], vertices[1], vertices[9], vertices[8]})
	p.AddFace([]*Vertex{vertices[1], vertices[16], vertices[17], vertices[3], vertices[13]})
	p.AddFace([]*Vertex{vertices[1], vertices[13], vertices[15], vertices[5], vertices[9]})
	p.AddFace([]*Vertex{vertices[2], vertices[12], vertices[14], vertices[6], vertices[10]})
	p.AddFace([]*Vertex{vertices[2], vertices[10], vertices[11], vertices[3], vertices[17]})
	p.AddFace([]*Vertex{vertices[3], vertices[11], vertices[7], vertices[15], vertices[13]})
	p.AddFace([]*Vertex{vertices[4], vertices[8], vertices[9], vertices[5], vertices[18]})
	p.AddFace([]*Vertex{vertices[4], vertices[18], vertices[19], vertices[6], vertices[14]})
	p.AddFace([]*Vertex{vertices[5], vertices[15], vertices[7], vertices[19], vertices[18]})
	p.AddFace([]*Vertex{vertices[6], vertices[19], vertices[7], vertices[11], vertices[10]})

	p.Normalize()

	return p
}

func Icosahedron() *Polyhedron {
	p := NewPolyhedron("Icosahedron")

	phi := (1.0 + math.Sqrt(goldenRatioBase)) / goldenRatioDivisor

	vertices := []*Vertex{
		p.AddVertex(Vector3{0, 1, phi}),
		p.AddVertex(Vector3{0, 1, -phi}),
		p.AddVertex(Vector3{0, -1, phi}),
		p.AddVertex(Vector3{0, -1, -phi}),

		p.AddVertex(Vector3{1, phi, 0}),
		p.AddVertex(Vector3{1, -phi, 0}),
		p.AddVertex(Vector3{-1, phi, 0}),
		p.AddVertex(Vector3{-1, -phi, 0}),

		p.AddVertex(Vector3{phi, 0, 1}),
		p.AddVertex(Vector3{phi, 0, -1}),
		p.AddVertex(Vector3{-phi, 0, 1}),
		p.AddVertex(Vector3{-phi, 0, -1}),
	}

	p.AddFace([]*Vertex{vertices[0], vertices[2], vertices[8]})
	p.AddFace([]*Vertex{vertices[0], vertices[8], vertices[4]})
	p.AddFace([]*Vertex{vertices[0], vertices[4], vertices[6]})
	p.AddFace([]*Vertex{vertices[0], vertices[6], vertices[10]})
	p.AddFace([]*Vertex{vertices[0], vertices[10], vertices[2]})

	p.AddFace([]*Vertex{vertices[3], vertices[1], vertices[9]})
	p.AddFace([]*Vertex{vertices[3], vertices[9], vertices[5]})
	p.AddFace([]*Vertex{vertices[3], vertices[5], vertices[7]})
	p.AddFace([]*Vertex{vertices[3], vertices[7], vertices[11]})
	p.AddFace([]*Vertex{vertices[3], vertices[11], vertices[1]})

	p.AddFace([]*Vertex{vertices[2], vertices[10], vertices[7]})
	p.AddFace([]*Vertex{vertices[2], vertices[7], vertices[5]})
	p.AddFace([]*Vertex{vertices[2], vertices[5], vertices[8]})

	p.AddFace([]*Vertex{vertices[8], vertices[5], vertices[9]})
	p.AddFace([]*Vertex{vertices[8], vertices[9], vertices[4]})

	p.AddFace([]*Vertex{vertices[4], vertices[9], vertices[1]})
	p.AddFace([]*Vertex{vertices[4], vertices[1], vertices[6]})

	p.AddFace([]*Vertex{vertices[6], vertices[1], vertices[11]})
	p.AddFace([]*Vertex{vertices[6], vertices[11], vertices[10]})

	p.AddFace([]*Vertex{vertices[10], vertices[11], vertices[7]})

	p.Normalize()

	return p
}

func GetSeed(symbol string) *Polyhedron {
	switch symbol {
	case "T":
		return Tetrahedron()
	case "C":
		return Cube()
	case "O":
		return Octahedron()
	case "D":
		return Dodecahedron()
	case "I":
		return Icosahedron()
	default:
		return nil
	}
}
