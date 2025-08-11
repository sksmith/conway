package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
)

func TestDualOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			original := test.poly()
			dual := conway.Dual(original)

			if !dual.IsValid() {
				t.Errorf("Dual of %s is not valid", test.name)
			}

			if dual.EulerCharacteristic() != 2 {
				t.Errorf("Dual of %s has wrong Euler characteristic: %d",
					test.name, dual.EulerCharacteristic())
			}

			doubleDual := conway.Dual(dual)

			if len(doubleDual.Vertices) != len(original.Vertices) {
				t.Errorf("Double dual vertex count mismatch: %d vs %d",
					len(doubleDual.Vertices), len(original.Vertices))
			}
		})
	}
}

func TestAmboOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			original := test.poly()
			ambo := conway.Ambo(original)

			if !ambo.IsValid() {
				t.Errorf("Ambo of %s is not valid", test.name)
			}

			if ambo.EulerCharacteristic() != 2 {
				t.Errorf("Ambo of %s has wrong Euler characteristic: %d",
					test.name, ambo.EulerCharacteristic())
			}

			if len(ambo.Vertices) != len(original.Edges) {
				t.Errorf("Ambo should have %d vertices (original edges), got %d",
					len(original.Edges), len(ambo.Vertices))
			}
		})
	}
}

func TestTruncateOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			original := test.poly()
			truncated := conway.Truncate(original)

			if !truncated.IsValid() {
				t.Errorf("Truncate of %s is not valid", test.name)
			}

			if truncated.EulerCharacteristic() != 2 {
				t.Errorf("Truncate of %s has wrong Euler characteristic: %d",
					test.name, truncated.EulerCharacteristic())
			}
		})
	}
}

func TestKisOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			original := test.poly()
			kis := conway.Kis(original)

			if !kis.IsValid() {
				t.Errorf("Kis of %s is not valid", test.name)
			}

			if kis.EulerCharacteristic() != 2 {
				t.Errorf("Kis of %s has wrong Euler characteristic: %d",
					test.name, kis.EulerCharacteristic())
			}

			expectedVertices := len(original.Vertices) + len(original.Faces)
			if len(kis.Vertices) != expectedVertices {
				t.Errorf("Kis should have %d vertices, got %d",
					expectedVertices, len(kis.Vertices))
			}
		})
	}
}

func TestJoinOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		poly func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			original := test.poly()
			joined := conway.Join(original)

			if !joined.IsValid() {
				t.Errorf("Join of %s is not valid", test.name)
			}

			if joined.EulerCharacteristic() != 2 {
				t.Errorf("Join of %s has wrong Euler characteristic: %d",
					test.name, joined.EulerCharacteristic())
			}
		})
	}
}

func TestCompoundOperations(t *testing.T) {
	t.Parallel()

	original := conway.Cube()

	t.Run("Ortho", func(t *testing.T) {
		t.Parallel()

		ortho := conway.Ortho(original)
		if !ortho.IsValid() {
			t.Error("Ortho operation produced invalid polyhedron")
		}
	})

	t.Run("Expand", func(t *testing.T) {
		t.Parallel()

		expand := conway.Expand(original)
		if !expand.IsValid() {
			t.Error("Expand operation produced invalid polyhedron")
		}
	})

	t.Run("Gyro", func(t *testing.T) {
		t.Parallel()

		gyro := conway.Gyro(original)
		if !gyro.IsValid() {
			t.Error("Gyro operation produced invalid polyhedron")
		}
	})

	t.Run("Snub", func(t *testing.T) {
		t.Parallel()

		snub := conway.Snub(original)
		if !snub.IsValid() {
			t.Error("Snub operation produced invalid polyhedron")
		}
	})
}

func TestOperationSymbols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		op       conway.Operation
		expected string
	}{
		{"Ambo", conway.AmboOp{}, "a"},
		{"Dual", conway.DualOp{}, "d"},
		{"Join", conway.JoinOp{}, "j"},
		{"Kis", conway.KisOp{}, "k"},
		{"Truncate", conway.TruncateOp{}, "t"},
		{"Ortho", conway.OrthoOp{}, "o"},
		{"Expand", conway.ExpandOp{}, "e"},
		{"Gyro", conway.GyroOp{}, "g"},
		{"Snub", conway.SnubOp{}, "s"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			symbol := test.op.Symbol()
			if symbol != test.expected {
				t.Errorf("Expected symbol %s for %s, got %s", test.expected, test.name, symbol)
			}
		})
	}
}

func TestOperationNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		op       conway.Operation
		expected string
	}{
		{"Ambo", conway.AmboOp{}, "ambo"},
		{"Dual", conway.DualOp{}, "dual"},
		{"Join", conway.JoinOp{}, "join"},
		{"Kis", conway.KisOp{}, "kis"},
		{"Truncate", conway.TruncateOp{}, "truncate"},
		{"Ortho", conway.OrthoOp{}, "ortho"},
		{"Expand", conway.ExpandOp{}, "expand"},
		{"Gyro", conway.GyroOp{}, "gyro"},
		{"Snub", conway.SnubOp{}, "snub"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			name := test.op.Name()
			if name != test.expected {
				t.Errorf("Expected name %s for operation, got %s", test.expected, name)
			}
		})
	}
}
