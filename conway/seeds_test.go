package conway

import "testing"

func TestSeedPolyhedra(t *testing.T) {
	tests := []struct {
		name string
		poly func() *Polyhedron
		V    int
		E    int
		F    int
	}{
		{"Tetrahedron", Tetrahedron, 4, 6, 4},
		{"Cube", Cube, 8, 12, 6},
		{"Octahedron", Octahedron, 6, 12, 8},
		{"Dodecahedron", Dodecahedron, 20, 30, 12},
		{"Icosahedron", Icosahedron, 12, 30, 20},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.poly()

			if len(p.Vertices) != test.V {
				t.Errorf("%s vertices: got %d, expected %d",
					test.name, len(p.Vertices), test.V)
			}

			if len(p.Edges) != test.E {
				t.Errorf("%s edges: got %d, expected %d",
					test.name, len(p.Edges), test.E)
			}

			if len(p.Faces) != test.F {
				t.Errorf("%s faces: got %d, expected %d",
					test.name, len(p.Faces), test.F)
			}

			if !p.IsValid() {
				t.Errorf("%s is not valid", test.name)
			}
		})
	}
}

func TestGetSeed(t *testing.T) {
	tests := []struct {
		symbol   string
		expected bool
	}{
		{"T", true},
		{"C", true},
		{"O", true},
		{"D", true},
		{"I", true},
		{"X", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.symbol, func(t *testing.T) {
			poly := GetSeed(test.symbol)
			if (poly != nil) != test.expected {
				t.Errorf("GetSeed(%s): got %v, expected existence %v",
					test.symbol, poly != nil, test.expected)
			}
		})
	}
}

func TestSeedValidity(t *testing.T) {
	seeds := []string{"T", "C", "O", "D", "I"}

	for _, symbol := range seeds {
		t.Run(symbol, func(t *testing.T) {
			p := GetSeed(symbol)
			if p == nil {
				t.Fatalf("GetSeed(%s) returned nil", symbol)
			}

			if !p.IsValid() {
				t.Errorf("Seed %s is not valid: %s", symbol, p.Stats())
			}

			if p.EulerCharacteristic() != 2 {
				t.Errorf("Seed %s has wrong Euler characteristic: %d",
					symbol, p.EulerCharacteristic())
			}
		})
	}
}

func TestSeedGeometry(t *testing.T) {
	t.Run("Tetrahedron", func(t *testing.T) {
		p := Tetrahedron()

		for _, f := range p.Faces {
			if f.Degree() != 3 {
				t.Errorf("Tetrahedron face should have 3 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 3 {
				t.Errorf("Tetrahedron vertex should have degree 3, got %d", v.Degree())
			}
		}
	})

	t.Run("Cube", func(t *testing.T) {
		p := Cube()

		for _, f := range p.Faces {
			if f.Degree() != 4 {
				t.Errorf("Cube face should have 4 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 3 {
				t.Errorf("Cube vertex should have degree 3, got %d", v.Degree())
			}
		}
	})

	t.Run("Octahedron", func(t *testing.T) {
		p := Octahedron()

		for _, f := range p.Faces {
			if f.Degree() != 3 {
				t.Errorf("Octahedron face should have 3 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 4 {
				t.Errorf("Octahedron vertex should have degree 4, got %d", v.Degree())
			}
		}
	})
}
