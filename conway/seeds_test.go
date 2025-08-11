package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
)

func TestSeedPolyhedra(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
		V    int
		E    int
		F    int
	}{
		{"conway.Tetrahedron", conway.Tetrahedron, 4, 6, 4},
		{"conway.Cube", conway.Cube, 8, 12, 6},
		{"conway.Octahedron", conway.Octahedron, 6, 12, 8},
		{"conway.Dodecahedron", conway.Dodecahedron, 20, 30, 12},
		{"conway.Icosahedron", conway.Icosahedron, 12, 30, 20},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

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
	t.Parallel()

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
			t.Parallel()

			poly := conway.GetSeed(test.symbol)
			if (poly != nil) != test.expected {
				t.Errorf("conway.GetSeed(%s): got %v, expected existence %v",
					test.symbol, poly != nil, test.expected)
			}
		})
	}
}

func TestSeedValidity(t *testing.T) {
	t.Parallel()

	seeds := []string{"T", "C", "O", "D", "I"}

	for _, symbol := range seeds {
		t.Run(symbol, func(t *testing.T) {
			t.Parallel()

			p := conway.GetSeed(symbol)
			if p == nil {
				t.Fatalf("conway.GetSeed(%s) returned nil", symbol)
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
	t.Parallel()

	t.Run("conway.Tetrahedron", func(t *testing.T) {
		t.Parallel()

		p := conway.Tetrahedron()

		for _, f := range p.Faces {
			if f.Degree() != 3 {
				t.Errorf("conway.Tetrahedron face should have 3 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 3 {
				t.Errorf("conway.Tetrahedron vertex should have degree 3, got %d", v.Degree())
			}
		}
	})

	t.Run("conway.Cube", func(t *testing.T) {
		t.Parallel()

		p := conway.Cube()

		for _, f := range p.Faces {
			if f.Degree() != 4 {
				t.Errorf("conway.Cube face should have 4 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 3 {
				t.Errorf("conway.Cube vertex should have degree 3, got %d", v.Degree())
			}
		}
	})

	t.Run("conway.Octahedron", func(t *testing.T) {
		t.Parallel()

		p := conway.Octahedron()

		for _, f := range p.Faces {
			if f.Degree() != 3 {
				t.Errorf("conway.Octahedron face should have 3 vertices, got %d", f.Degree())
			}
		}

		for _, v := range p.Vertices {
			if v.Degree() != 4 {
				t.Errorf("conway.Octahedron vertex should have degree 4, got %d", v.Degree())
			}
		}
	})
}
