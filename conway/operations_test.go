package conway

import "testing"

func TestDualOperation(t *testing.T) {
	tests := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.poly()
			dual := Dual(original)

			if !dual.IsValid() {
				t.Errorf("Dual of %s is not valid", test.name)
			}

			if dual.EulerCharacteristic() != 2 {
				t.Errorf("Dual of %s has wrong Euler characteristic: %d",
					test.name, dual.EulerCharacteristic())
			}

			doubleDual := Dual(dual)

			if len(doubleDual.Vertices) != len(original.Vertices) {
				t.Errorf("Double dual vertex count mismatch: %d vs %d",
					len(doubleDual.Vertices), len(original.Vertices))
			}
		})
	}
}

func TestAmboOperation(t *testing.T) {
	tests := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.poly()
			ambo := Ambo(original)

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
	tests := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.poly()
			truncated := Truncate(original)

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
	tests := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.poly()
			kis := Kis(original)

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
	tests := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.poly()
			joined := Join(original)

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
	original := Cube()

	t.Run("Ortho", func(t *testing.T) {
		ortho := Ortho(original)
		if !ortho.IsValid() {
			t.Error("Ortho operation produced invalid polyhedron")
		}
	})

	t.Run("Expand", func(t *testing.T) {
		expand := Expand(original)
		if !expand.IsValid() {
			t.Error("Expand operation produced invalid polyhedron")
		}
	})

	t.Run("Gyro", func(t *testing.T) {
		gyro := Gyro(original)
		if !gyro.IsValid() {
			t.Error("Gyro operation produced invalid polyhedron")
		}
	})

	t.Run("Snub", func(t *testing.T) {
		snub := Snub(original)
		if !snub.IsValid() {
			t.Error("Snub operation produced invalid polyhedron")
		}
	})
}

func TestOperationSymbols(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected string
	}{
		{"Ambo", AmboOp{}, "a"},
		{"Dual", DualOp{}, "d"},
		{"Join", JoinOp{}, "j"},
		{"Kis", KisOp{}, "k"},
		{"Truncate", TruncateOp{}, "t"},
		{"Ortho", OrthoOp{}, "o"},
		{"Expand", ExpandOp{}, "e"},
		{"Gyro", GyroOp{}, "g"},
		{"Snub", SnubOp{}, "s"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			symbol := test.op.Symbol()
			if symbol != test.expected {
				t.Errorf("Expected symbol %s for %s, got %s", test.expected, test.name, symbol)
			}
		})
	}
}

func TestOperationNames(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		expected string
	}{
		{"Ambo", AmboOp{}, "ambo"},
		{"Dual", DualOp{}, "dual"},
		{"Join", JoinOp{}, "join"},
		{"Kis", KisOp{}, "kis"},
		{"Truncate", TruncateOp{}, "truncate"},
		{"Ortho", OrthoOp{}, "ortho"},
		{"Expand", ExpandOp{}, "expand"},
		{"Gyro", GyroOp{}, "gyro"},
		{"Snub", SnubOp{}, "snub"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name := test.op.Name()
			if name != test.expected {
				t.Errorf("Expected name %s for operation, got %s", test.expected, name)
			}
		})
	}
}
