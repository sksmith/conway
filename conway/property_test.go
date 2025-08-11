package conway

import (
	"testing"
)

// TestDualInvolution tests the dual involution property: dd(P) = P
// This is a fundamental property of the dual operation
func TestDualInvolution(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
		{"Dodecahedron", Dodecahedron},
		{"Icosahedron", Icosahedron},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.poly()
			dual1 := Dual(original)
			dual2 := Dual(dual1)

			// Check that dd(P) has same topology as P
			if len(dual2.Vertices) != len(original.Vertices) {
				t.Errorf("Dual involution failed: vertex count mismatch. Original: %d, dd(P): %d",
					len(original.Vertices), len(dual2.Vertices))
			}

			if len(dual2.Edges) != len(original.Edges) {
				t.Errorf("Dual involution failed: edge count mismatch. Original: %d, dd(P): %d",
					len(original.Edges), len(dual2.Edges))
			}

			if len(dual2.Faces) != len(original.Faces) {
				t.Errorf("Dual involution failed: face count mismatch. Original: %d, dd(P): %d",
					len(original.Faces), len(dual2.Faces))
			}

			// Check Euler characteristic preservation
			if dual2.EulerCharacteristic() != original.EulerCharacteristic() {
				t.Errorf("Dual involution failed: Euler characteristic mismatch. Original: %d, dd(P): %d",
					original.EulerCharacteristic(), dual2.EulerCharacteristic())
			}

			// Validate that both are valid polyhedra
			if !original.IsValid() {
				t.Errorf("Original %s is not valid", tc.name)
			}
			if !dual1.IsValid() {
				t.Errorf("First dual of %s is not valid", tc.name)
			}
			if !dual2.IsValid() {
				t.Errorf("Second dual of %s is not valid", tc.name)
			}
		})
	}
}

// TestEulerCharacteristicPreservation tests that all operations preserve Euler characteristic
func TestEulerCharacteristicPreservation(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	operations := []struct {
		name string
		op   func(*Polyhedron) *Polyhedron
	}{
		{"Dual", Dual},
		{"Ambo", Ambo},
		{"Truncate", Truncate},
		{"Kis", Kis},
		{"Join", Join},
	}

	for _, tc := range testCases {
		for _, op := range operations {
			t.Run(tc.name+"_"+op.name, func(t *testing.T) {
				original := tc.poly()
				result := op.op(original)

				originalEuler := original.EulerCharacteristic()
				resultEuler := result.EulerCharacteristic()

				if originalEuler != resultEuler {
					t.Errorf("Operation %s on %s failed to preserve Euler characteristic. Original: %d, Result: %d",
						op.name, tc.name, originalEuler, resultEuler)
				}

				if originalEuler != 2 {
					t.Errorf("Original %s has incorrect Euler characteristic: %d (expected 2)", tc.name, originalEuler)
				}

				if resultEuler != 2 {
					t.Errorf("Result %s(%s) has incorrect Euler characteristic: %d (expected 2)", op.name, tc.name, resultEuler)
				}
			})
		}
	}
}

// TestOperationValidation tests that all operations produce valid polyhedra
func TestOperationValidation(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	operations := []struct {
		name string
		op   func(*Polyhedron) *Polyhedron
	}{
		{"Dual", Dual},
		{"Ambo", Ambo},
		{"Truncate", Truncate},
		{"Kis", Kis},
		{"Join", Join},
	}

	for _, tc := range testCases {
		for _, op := range operations {
			t.Run(tc.name+"_"+op.name, func(t *testing.T) {
				original := tc.poly()
				result := op.op(original)

				if !result.IsValid() {
					t.Errorf("Operation %s on %s produced invalid polyhedron", op.name, tc.name)
				}

				// Test comprehensive validation
				if err := result.ValidateComplete(); err != nil {
					t.Errorf("Operation %s on %s failed comprehensive validation: %v", op.name, tc.name, err)
				}
			})
		}
	}
}

// TestDualVertexFaceCorrespondence tests the vertex-face correspondence in dual operation
func TestDualVertexFaceCorrespondence(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.poly()
			dual := Dual(original)

			// In the dual operation:
			// - Each vertex of the original becomes a face of the dual
			// - Each face of the original becomes a vertex of the dual
			if len(original.Vertices) != len(dual.Faces) {
				t.Errorf("Dual operation failed vertex-face correspondence. Original vertices: %d, Dual faces: %d",
					len(original.Vertices), len(dual.Faces))
			}

			if len(original.Faces) != len(dual.Vertices) {
				t.Errorf("Dual operation failed face-vertex correspondence. Original faces: %d, Dual vertices: %d",
					len(original.Faces), len(dual.Vertices))
			}

			// Edge count should remain the same
			if len(original.Edges) != len(dual.Edges) {
				t.Errorf("Dual operation failed edge preservation. Original edges: %d, Dual edges: %d",
					len(original.Edges), len(dual.Edges))
			}
		})
	}
}

// TestAmboDoubleProperty tests that ambo applied twice (expand) produces predictable results
func TestAmboDoubleProperty(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.poly()
			ambo1 := Ambo(original)
			expand := Ambo(ambo1) // aa = expand

			// Validate intermediate and final results
			if !ambo1.IsValid() {
				t.Errorf("First ambo of %s is invalid", tc.name)
			}

			if !expand.IsValid() {
				t.Errorf("Expand of %s is invalid", tc.name)
			}

			// Check Euler characteristic preservation
			if expand.EulerCharacteristic() != 2 {
				t.Errorf("Expand of %s has incorrect Euler characteristic: %d", tc.name, expand.EulerCharacteristic())
			}
		})
	}
}

// TestOperationComposition tests that composition of operations works correctly
func TestOperationComposition(t *testing.T) {
	cube := Cube()

	// Test dtC (dual of truncated cube)
	truncated := Truncate(cube)
	dtC := Dual(truncated)

	if !truncated.IsValid() {
		t.Error("Truncated cube is invalid")
	}

	if !dtC.IsValid() {
		t.Error("Dual of truncated cube is invalid")
	}

	if dtC.EulerCharacteristic() != 2 {
		t.Errorf("dtC has incorrect Euler characteristic: %d", dtC.EulerCharacteristic())
	}

	// Test comprehensive validation
	if err := dtC.ValidateComplete(); err != nil {
		t.Errorf("dtC failed comprehensive validation: %v", err)
	}
}

// TestTopologyConsistency tests that topology remains consistent throughout operations
func TestTopologyConsistency(t *testing.T) {
	testCases := []struct {
		name string
		poly func() *Polyhedron
	}{
		{"Tetrahedron", Tetrahedron},
		{"Cube", Cube},
		{"Octahedron", Octahedron},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.poly()

			// Check that each edge connects exactly 2 vertices
			for _, edge := range original.Edges {
				if edge.V1 == nil || edge.V2 == nil {
					t.Errorf("%s has edge with nil vertex", tc.name)
				}
				if edge.V1.ID == edge.V2.ID {
					t.Errorf("%s has edge connecting vertex to itself", tc.name)
				}
			}

			// Check that each face has at least 3 vertices
			for _, face := range original.Faces {
				if len(face.Vertices) < 3 {
					t.Errorf("%s has face with < 3 vertices", tc.name)
				}
				if len(face.Edges) < 3 {
					t.Errorf("%s has face with < 3 edges", tc.name)
				}
			}

			// Check vertex-edge consistency
			for _, vertex := range original.Vertices {
				for _, edge := range vertex.Edges {
					if edge.V1.ID != vertex.ID && edge.V2.ID != vertex.ID {
						t.Errorf("%s has inconsistent vertex-edge relationship", tc.name)
					}
				}
			}
		})
	}
}
